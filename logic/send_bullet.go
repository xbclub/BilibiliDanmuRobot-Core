package logic

import (
	"context"
	"github.com/xbclub/BilibiliDanmuRobot-Core/biliapi"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var sender = &BulletSender{
	bulletChan: make(chan entity.Bullet, 1000),
}

// sendInterval 中央发送队列的全局出口限速：B 站弹幕实际约 3 秒一条，留 0.5s 宽限
const sendInterval = 3500 * time.Millisecond

type BulletSender struct {
	bulletChan chan entity.Bullet
}

func PushToBulletSender(msg string, reply ...*entity.DanmuMsgTextReplyInfo) {
	logx.Info("PushToBulletSender成功", msg)
	if sender == nil {
		return
	}
	bullet := entity.Bullet{
		Msg:   msg,
		Reply: reply,
	}
	sender.bulletChan <- bullet
}

func StartSendBullet(ctx context.Context, svcCtx *svc.ServiceContext) {

	for {
		select {
		case <-ctx.Done():
			goto END
		case bullet := <-sender.bulletChan:
			svcCtx.EventRLock()
			danmuLen := svcCtx.Config.DanmuLen
			svcCtx.EventRUnlock()
			for _, msgs := range SplitMessage(bullet.Msg, danmuLen) {
				if err := biliapi.Send(msgs, svcCtx, bullet.Reply...); err != nil {
					logx.Errorf("弹幕发送失败：%s msg: %s", err, msgs)
				} else {
					logx.Infof("弹幕发送成功：%s", msgs)
				}
				//fmt.Println(msgs)
				time.Sleep(sendInterval) // 防止弹幕发送过快
			}
		}

	}
END:
}
