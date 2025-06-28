package api

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	log "github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RoomInfo
// api https://api.live.bilibili.com/room/v1/Room/room_init?id={} response
type RoomInfo struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		RoomId      int   `json:"room_id"`
		ShortId     int   `json:"short_id"`
		Uid         int   `json:"uid"`
		NeedP2P     int   `json:"need_p2p"`
		IsHidden    bool  `json:"is_hidden"`
		IsLocked    bool  `json:"is_locked"`
		IsPortrait  bool  `json:"is_portrait"`
		LiveStatus  int   `json:"live_status"`
		HiddenTill  int   `json:"hidden_till"`
		LockTill    int   `json:"lock_till"`
		Encrypted   bool  `json:"encrypted"`
		PwdVerified bool  `json:"pwd_verified"`
		LiveTime    int64 `json:"live_time"`
		RoomShield  int   `json:"room_shield"`
		IsSp        int   `json:"is_sp"`
		SpecialType int   `json:"special_type"`
	} `json:"data"`
}

// DanmuInfo
// api https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id={}&type=0 response
type DanmuInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Group            string  `json:"group"`
		BusinessId       int     `json:"business_id"`
		RefreshRowFactor float64 `json:"refresh_row_factor"`
		RefreshRate      int     `json:"refresh_rate"`
		MaxDelay         int     `json:"max_delay"`
		Token            string  `json:"token"`
		HostList         []struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			WssPort int    `json:"wss_port"`
			WsPort  int    `json:"ws_port"`
		} `json:"host_list"`
	} `json:"data"`
}

func GetUid(cookie string) (int, string, error) {
	headers := &http.Header{}
	headers.Set("cookie", cookie)
	resp, err := HttpGet("https://api.bilibili.com/x/web-interface/nav", headers)
	if err != nil {
		return 0, "", err
	}
	j := gjson.ParseBytes(resp)
	if j.Get("code").Int() != 0 || !j.Get("data.isLogin").Bool() {
		return 0, "", errors.New(j.Get("message").String())
	}
	if !j.Get("data.wbi_img").Exists() ||
		!j.Get("data.wbi_img.img_url").Exists() ||
		!j.Get("data.wbi_img.sub_url").Exists() {
		return 0, "", errors.New("wbi_img not found in response")
	}
	imgUrl := j.Get("data.wbi_img.img_url").String()
	subUrl := j.Get("data.wbi_img.sub_url").String()
	// 定义正则表达式，提取中间的32位十六进制字符串
	re := regexp.MustCompile(`/(\w{32})\.png`)
	// 匹配 img_url
	imgMatches := re.FindStringSubmatch(imgUrl)
	if len(imgMatches) < 2 {
		panic("img_url 中未找到32位字符串")
	}
	wbiImgBa := imgMatches[1]

	// 匹配 sub_url
	subMatches := re.FindStringSubmatch(subUrl)
	if len(subMatches) < 2 {
		panic("sub_url 中未找到32位字符串")
	}
	wbiSubBa := subMatches[1]

	// 拼接两个32位字符串为64位
	wbiImgSub := wbiImgBa + wbiSubBa
	if len(wbiImgSub) != 64 {
		panic("拼接后字符串长度不为64")
	}

	// 索引重排表
	mixinKeyEncTab := []int{
		46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
		33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
		61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
		36, 20, 34, 44, 52,
	}
	// 使用索引重排表对 wbiImgSub 进行重排
	var newString []byte
	for _, idx := range mixinKeyEncTab {
		if idx < 0 || idx >= len(wbiImgSub) {
			panic("索引越界: " + fmt.Sprintf("%d", idx))
		}
		newString = append(newString, wbiImgSub[idx])
	}

	// 截取前32位作为最终密钥
	wbiMixinKey := string(newString[:32])
	return int(j.Get("data.mid").Int()), wbiMixinKey, nil
}

// 使用 url.Values 并正确编码
func toWbiParamSafe(params string, wbiMixinKey string) string {
	values, _ := url.ParseQuery(params)

	if _, exists := values["wts"]; !exists {
		values.Set("wts", strconv.FormatInt(time.Now().Unix(), 10))
	}

	// 按照 key 字典序排序
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接 key=value 形式
	var orderedPairs []string
	for _, k := range keys {
		orderedPairs = append(orderedPairs, k+"="+values.Get(k))
	}
	valuesEncoded := strings.Join(orderedPairs, "&")

	// 计算签名
	hash := md5.Sum([]byte(valuesEncoded + wbiMixinKey))
	md5Hash := fmt.Sprintf("%x", hash)

	return valuesEncoded + "&w_rid=" + md5Hash
}

type BuvidData struct {
	Code int `json:"code"`
	Data struct {
		B3 string `json:"b_3"`
		B4 string `json:"b_4"`
	} `json:"data"`
	Message string `json:"message"`
}

func GetBuvid3A4() (buvid3, buvid4 string, err error) {
	headers := &http.Header{}
	roomIDurl := "https://api.bilibili.com/x/frontend/finger/spi"
	result := &BuvidData{}
	err = GetJsonWithHeader(roomIDurl, headers, result)
	if err != nil {
		return "", "", err
	}
	if result.Code != 0 {
		return "", "", errors.New(result.Message)
	}
	return result.Data.B3, result.Data.B4, nil
}
func GetDanmuInfo(roomID int, cookie string, wbiMixinKey string) (*DanmuInfo, error) {
	result := &DanmuInfo{}
	headers := &http.Header{}
	headers.Set("cookie", cookie)
	//headers.Set("content-type", "application/json")
	//headers.Set("accept-encoding", "gzip, deflate, br, zstd")
	//headers.Set("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	params := fmt.Sprintf("id=%d&type=0", roomID)
	params = toWbiParamSafe(params, wbiMixinKey)
	roomIDurl := fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?%s", params)
	err := GetJsonWithHeader(roomIDurl, headers, result)
	if err != nil {
		return nil, err
	}
	if result.Code == -352 {
		log.Errorf("request data: %+v", roomIDurl)
		log.Errorf("header data: %+v", headers)
		log.Errorf("response data: %+v", result)
		return nil, errors.New("触发风控")
	}
	return result, nil
}

func GetRoomInfo(roomID int) (*RoomInfo, error) {
	result := &RoomInfo{}
	err := GetJson(fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/room_init?id=%d", roomID), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetRoomRealID(roomID int) (string, error) {
	res, err := GetRoomInfo(roomID)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(res.Data.RoomId), nil
}
