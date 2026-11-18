package pk

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xbclub/BilibiliDanmuRobot-Core/biliapi"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var pkGiver = &PKGiver{
	pkFilter: map[int]time.Time{},
	locked:   new(sync.Mutex),
	tableMu:  sync.RWMutex{},
	pkChan:   make(chan *int, 1000),
}

type PKGiver struct {
	pkFilter map[int]time.Time
	locked   *sync.Mutex
	tableMu  sync.RWMutex
	pkChan   chan *int
}

func PushToPKChan(g *int) {
	pkGiver.pkChan <- g
}

// Run 启动 PK 信息 worker，同一房间 10 秒内只播报一次
func Run(ctx context.Context, svcCtx *svc.ServiceContext) {

	var g *int
	var w = 10 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			if len(pkGiver.pkFilter) > 0 {
				pkGiver.locked.Lock()
				for k, v := range pkGiver.pkFilter {
					if v.Add(w).Unix() < time.Now().Unix() {
						delete(pkGiver.pkFilter, k)
						logx.Debugf("pk room %v 已从重复过滤列表移除", k)
					}
				}
				pkGiver.locked.Unlock()
			}

			t.Reset(w)
		case g = <-pkGiver.pkChan:
			pkGiver.locked.Lock()
			if value, ok := pkGiver.pkFilter[*g]; ok && value.Add(w).Unix() >= time.Now().Unix() {
				logx.Debugf("pk room %v 10秒内重复获取数据已被过滤", *g)
			} else {
				logx.Debugf("正在处理pk信息")
				handlerPK(*g, svcCtx)
				logx.Debug(*g)
				pkGiver.pkFilter[*g] = time.Now()
			}
			pkGiver.locked.Unlock()
			logx.Debugf("pk room %v 已进入重复过滤列表", *g)
		}
	}
END:
}
func handlerPK(roomid int, svcCtx *svc.ServiceContext) {

	toplist := &entity.TopListInfo{}
	toplistalive := 0
	rankcount := 0

	userinfo, err := biliapi.Userinfo(roomid)
	if err != nil {
		logx.Error(err)
		return
	}
	toppage := 1
	listInfo, err := biliapi.TopListInfo(roomid, userinfo.Data.Info.Uid, toppage)
	if err != nil {
		logx.Error(err)
		logic.PushToBulletSender("PK信息获取失败!")
		return
	}

	tmpPage := listInfo.Data.Info.Page
	for toppage += 1; toppage <= tmpPage; toppage++ {
		toplist, err = biliapi.TopListInfo(roomid, userinfo.Data.Info.Uid, toppage)
		if err != nil {
			logx.Error(err)
			continue
		}
		listInfo.Data.List = append(listInfo.Data.List, toplist.Data.List...)
	}

	svcCtx.OtherSideUid.Clear()

	for _, data := range listInfo.Data.List {
		svcCtx.OtherSideUid.Add(data.Uid)
	}
	rankListInfo, err := biliapi.RankListInfo(roomid, userinfo.Data.Info.Uid, 1)
	if err != nil {
		logic.PushToBulletSender("PK信息获取失败!")
		logx.Error(err)
		return
	}
	for _, data := range rankListInfo.Data.OnlineRankItem {
		rankcount += data.Score
		if data.GuardLevel > 0 {
			toplistalive++
		}
		svcCtx.OtherSideUid.Add(data.Uid)
	}

	// 获取剩下的高能榜
	if len(rankListInfo.Data.OnlineRankItem) > 0 && len(rankListInfo.Data.OnlineRankItem) < rankListInfo.Data.OnlineNum {
		totalPage := rankListInfo.Data.OnlineNum / len(rankListInfo.Data.OnlineRankItem)
		if rankListInfo.Data.OnlineNum%50 > 0 {
			totalPage++
		}

		p := 2
		for ; p <= totalPage; p++ {
			rankListInfoTmp, err := biliapi.RankListInfo(roomid, userinfo.Data.Info.Uid, 1)
			if err != nil {
				logx.Error(err)
				continue
			}
			rankListInfo.Data.OnlineRankItem = append(rankListInfo.Data.OnlineRankItem, rankListInfoTmp.Data.OnlineRankItem...)
		}
		for _, data := range rankListInfo.Data.OnlineRankItem {
			if data.GuardLevel > 0 {
				toplistalive++
			}
			svcCtx.OtherSideUid.Add(data.Uid)
		}
	}

	logic.PushToBulletSender(fmt.Sprintf("当前对手:%v", userinfo.Data.Info.Uname))
	logic.PushToBulletSender(fmt.Sprintf("共%v船，%v粉", listInfo.Data.Info.Num, userinfo.Data.FollowerNum))
	logic.PushToBulletSender(fmt.Sprintf("当前%v船在线，高能榜%v人", toplistalive, rankListInfo.Data.OnlineNum))
	logic.PushToBulletSender(fmt.Sprintf("榜前50贡献%v分", rankcount))
}
