package client

import (
	"context"
	"errors"
	"fmt"
	log "github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/api"
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/packet"
	"github.com/xbclub/BilibiliDanmuRobot-Core/browserclient"
	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
)

const (
	// defaultDanmuHost 拿不到 getDanmuInfo（游客触发 -352 风控等）时的降级弹幕服务器。
	// 该服务器不校验进房 key，因此配合空 token 即可连接。
	defaultDanmuHost = "broadcastlv.chat.bilibili.com"

	// writeTimeout 单次写操作的超时。gorilla 的写没有默认截止时间，
	// 对端不读时会一直占着锁，进而拖住 Stop()；加一个上界即可。
	writeTimeout = 10 * time.Second
)

// gorillaManagedHeaders 由 gorilla/websocket 自己写入的握手头。
// 这些头不允许出现在 Dial 的 requestHeader 里，否则 gorilla 直接返回
// "websocket: duplicate header not allowed"。
// 键使用 http.CanonicalHeaderKey 的形式。
var gorillaManagedHeaders = map[string]struct{}{
	"Upgrade":                  {},
	"Connection":               {},
	"Sec-Websocket-Key":        {},
	"Sec-Websocket-Version":    {},
	"Sec-Websocket-Extensions": {},
	"Sec-Websocket-Protocol":   {},
}

type Client struct {
	conn                *websocket.Conn
	RoomID              int
	Uid                 int
	Buvid               string
	Cookie              string
	WbiMixinKey         string
	token               string
	host                string
	hostList            []string
	retryCount          int
	eventHandlers       *eventHandlers
	customEventHandlers *customEventHandlers
	cancel              context.CancelFunc
	done                <-chan struct{}
	lock                sync.RWMutex
}

// NewClient 创建一个新的弹幕 client
func NewClient(roomID int) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		RoomID:              roomID,
		retryCount:          0,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		done:                ctx.Done(),
		cancel:              cancel,
		lock:                sync.RWMutex{},
	}
}

func (c *Client) SetCookie(cookie string) {
	c.Cookie = cookie
}

// buvidFromCookie 从 cookie 串中提取进房包 buvid 字段的值。
// 优先取真实 buvid3（handler 层 ensureBuvid 已把 buvid3/buvid4 拼进 cookie），
// 取不到再回退 _uuid；同时兼容「值后跟分号」与「值在串尾」两种写法。
func buvidFromCookie(cookie string) string {
	for _, key := range []string{"buvid3", "_uuid"} {
		re := regexp.MustCompile(key + `=([^;]+);?`)
		if m := re.FindStringSubmatch(cookie); len(m) > 1 && m[1] != "" {
			return m[1]
		}
	}
	return ""
}

// init 初始化 获取真实 RoomID 和 弹幕服务器 host
func (c *Client) init() error {
	if c.Cookie != "" {
		if !strings.Contains(c.Cookie, "bili_jct") || !strings.Contains(c.Cookie, "SESSDATA") {
			log.Errorf("cannot found account token")
			return errors.New("账号未登录")
		}
		uid, wbiMixinKey, err := api.GetUid(c.Cookie)
		if err != nil {
			return errors.New(fmt.Sprintf("get uid error: %v", err))
		}
		c.Uid = uid
		c.WbiMixinKey = wbiMixinKey
		// 进房包 buvid 字段应传真实 buvid3，cookie 里没有时再回退 _uuid
		if v := buvidFromCookie(c.Cookie); v != "" {
			c.Buvid = v
		}
		b3, _, err := api.GetBuvid3A4()
		if err != nil {
			return err
		}
		// 按浏览器格式拼接（"k=v; k=v"），不再依赖上游 cookie 以分号结尾
		if len(c.Cookie) > 0 && !strings.HasSuffix(c.Cookie, ";") {
			c.Cookie += "; "
		}
		c.Cookie += "buvid3=" + b3
	}

	roomInfo, err := api.GetRoomInfo(c.RoomID)
	// 失败降级
	if err != nil || roomInfo.Code != 0 {
		return errors.New(fmt.Sprintf("room=%d init GetRoomInfo fialed, %s", c.RoomID, err))
	}
	c.RoomID = roomInfo.Data.RoomId
	if c.host == "" {
		info, err := api.GetDanmuInfo(c.RoomID, c.Cookie, c.WbiMixinKey)
		// 游客 / 风控（-352）时 GetDanmuInfo 返回 (nil, err)。
		// 此时只能走默认 host + 空 token，绝不能解引用 info.Data（否则 nil panic）。
		if err != nil || info == nil {
			log.Errorf("get DanmuInfo failed, fallback to default host %s, err: %v", defaultDanmuHost, err)
			c.hostList = []string{defaultDanmuHost}
			c.token = ""
		} else {
			if info.Data.Token != "" {
				c.token = info.Data.Token
			}
			for _, h := range info.Data.HostList {
				if h.Host != "" {
					c.hostList = append(c.hostList, h.Host)
				}
			}
			if len(c.hostList) == 0 {
				log.Errorf("DanmuInfo has empty host list, fallback to default host %s", defaultDanmuHost)
				c.hostList = []string{defaultDanmuHost}
			}
		}
	}
	return nil
}

// wsDialer TLS 握手走 Chrome 指纹（browserclient），HTTP/1.1 upgrade 语义不变
var wsDialer = websocket.Dialer{
	NetDialTLSContext: browserclient.DialTLSContext,
	HandshakeTimeout:  45 * time.Second,
}

// wsHeaders 返回 WS 握手用的完整浏览器特征头。
//
// 取值统一来自 webprofile.WSHeaders()（UA / Origin / Accept-Encoding /
// Accept-Language / Pragma / Cache-Control / Sec-Fetch-* / sec-ch-ua*），
// 与 API 侧保持同一套浏览器身份。
//
// Upgrade / Connection / Sec-WebSocket-* 由 gorilla 自己在握手时写入，
// 这里必须过滤掉：gorilla 遇到重复的这几项会直接返回错误。
func wsHeaders() http.Header {
	header := http.Header{}
	for k, v := range webprofile.WSHeaders() {
		if _, managed := gorillaManagedHeaders[http.CanonicalHeaderKey(k)]; managed {
			continue
		}
		// 用 Set 而不是直接赋值：Set 会规范化键名，
		// 让上面的 managed 判断对任意大小写写法都成立。
		header.Set(k, v)
	}
	return header
}

func (c *Client) connect() error {
	reqHeader := wsHeaders()
retry:
	c.host = c.hostList[c.retryCount%len(c.hostList)]
	c.retryCount++
	conn, res, err := wsDialer.Dial(fmt.Sprintf("wss://%s/sub", c.host), reqHeader)
	if err != nil {
		log.Errorf("connect dial failed, retry %d times", c.retryCount)
		time.Sleep(2 * time.Second)
		goto retry
	}
	c.conn = conn
	_ = res.Body.Close()
	if err = c.sendEnterPacket(); err != nil {
		log.Errorf("failed to send enter packet, retry %d times", c.retryCount)
		time.Sleep(2 * time.Second)
		goto retry
	}
	return nil
}

func (c *Client) wsLoop() {
	for {
		select {
		case <-c.done:
			return
		default:
			msgType, data, err := c.conn.ReadMessage()
			if err != nil {
				log.Error("ws message read failed, reconnecting")
				time.Sleep(time.Duration(3) * time.Millisecond)
				_ = c.connect()
				continue
			}
			if msgType != websocket.BinaryMessage {
				log.Error("packet not binary")
				continue
			}
			for _, pkt := range packet.DecodePacket(data).Parse() {
				go c.Handle(pkt)
			}
		}
	}
}

func (c *Client) heartBeatLoop() {
	pkt := packet.NewHeartBeatPacket()
	for {
		select {
		case <-c.done:
			return
		case <-time.After(30 * time.Second):
			c.lock.Lock()
			// gorilla 的写没有默认截止时间，对端不读时会一直占着锁；
			// 写前设置上界，避免拖住 Stop()。
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
				log.Error(err)
			}
			c.lock.Unlock()
			log.Debug("send: HeartBeat")
		}
	}
}

// Start 启动弹幕 Client 初始化并连接 ws、发送心跳包
func (c *Client) Start() error {
	if err := c.init(); err != nil {
		return err
	}
	if err := c.connect(); err != nil {
		return err
	}
	go c.wsLoop()
	go c.heartBeatLoop()
	return nil
}

// Stop 停止弹幕 Client
func (c *Client) Stop() {
	c.cancel()
}

func (c *Client) SetHost(host string) {
	c.host = host
}

// UseDefaultHost 使用默认 host broadcastlv.chat.bilibili.com
func (c *Client) UseDefaultHost() {
	c.hostList = []string{defaultDanmuHost}
}

func (c *Client) sendEnterPacket() error {
	pkt := packet.NewEnterPacket(c.Uid, c.Buvid, c.RoomID, c.token)
	c.lock.Lock()
	defer c.lock.Unlock()
	_ = c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		return err
	}
	log.Debugf("send: EnterPacket")
	return nil
}
