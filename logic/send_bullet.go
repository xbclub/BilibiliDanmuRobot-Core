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

type BulletSender struct {
	bulletChan chan entity.Bullet
}

func PushToBulletSender(msg string, reply ...*entity.DanmuMsgTextReplyInfo) {
	logx.Info("PushToBulletSender成功", msg)
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
			for _, msgs := range SplitMessage(bullet.Msg, svcCtx.Config.DanmuLen) {
				if err := biliapi.Send(msgs, svcCtx, bullet.Reply...); err != nil {
					logx.Errorf("弹幕发送失败：%s msg: %s", err, msgs)
				} else {
					logx.Infof("弹幕发送成功：%s", msgs)
				}
				//fmt.Println(msgs)
				time.Sleep(1 * time.Second) // 防止弹幕发送过快
			}
		}

	}
END:
}
