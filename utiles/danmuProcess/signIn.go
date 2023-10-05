package danmuProcess

import (
	"database/sql"
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	_ "github.com/glebarez/go-sqlite"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"time"
)

type SignIn struct {
	danmuContent *string
	fromUser     *message.User
	svcCtx       *svc.ServiceContext
}

func (signIn *SignIn) DoDanmuProcess() {
	if *signIn.danmuContent != "签到" {
		return
	}

	dbDirectory := "./db"
	info, err := os.Stat(dbDirectory)
	if os.IsNotExist(err) || !info.IsDir() {
		os.MkdirAll(dbDirectory, 0777)
	}
	dbFile := fmt.Sprintf("%s/%s?_pragma=busy_timeout(5000)", dbDirectory, "signin.db")

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		logx.Errorf("数据库打开失败：%s", err)
		return
	}

	sqlCommandBase := `
		drop table if exists '%d';
		create table '%d'(uid,lastday,numberOfConsecutiveCheckInDays)`

	sqlCommand := fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, signIn.svcCtx.Config.RoomId)
	if _, err = db.Exec(sqlCommand); err != nil {
		db.Close()
		logx.Errorf("SQL执行失败：%s", err)
		return
	}

	sqlCommandBase = "select * from '%d' where uid=%d;"
	sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, signIn.fromUser.Uid)
	rows, err := db.Query(sqlCommand)
	if err != nil {
		db.Close()
		logx.Errorf("SQL执行失败：%s", err)
		return
	}

	// 获取当前时间
	now := time.Now().In(time.Local)
	if rows.Next() {
		sqlCommandBase = "update '%d' set lastday=%d, numberOfConsecutiveCheckInDays=%d where uid=%d;"
		var uid, lastday int
		var numberOfConsecutiveCheckInDays int64
		rows.Scan(uid, lastday, numberOfConsecutiveCheckInDays)

		// 将时间戳转换为时间对象（中国时区）
		lastTime := time.Unix(numberOfConsecutiveCheckInDays, 0).In(time.Local)
		if now.Year() == lastTime.Year() && now.Month() == lastTime.Month() && now.Day() == lastTime.Day()+1 {
			sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, now.Unix(), numberOfConsecutiveCheckInDays+1, signIn.fromUser.Uid)
			if _, err = db.Exec(sqlCommand); err != nil {
				db.Close()
				logx.Errorf("SQL执行失败：%s", err)
				return
			}
			strMessage := fmt.Sprintf("%s,连续签到第%d天", signIn.fromUser.Uname, numberOfConsecutiveCheckInDays+1)
			logic.PushToBulletSender(strMessage)
		} else if now.Year() == lastTime.Year() && now.Month() == lastTime.Month() && now.Day() != lastTime.Day() {
			sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, now.Unix(), numberOfConsecutiveCheckInDays+1, signIn.fromUser.Uid)
			if _, err = db.Exec(sqlCommand); err != nil {
				db.Close()
				logx.Errorf("SQL执行失败：%s", err)
				return
			}
			strMessage := fmt.Sprintf("%s,连续签到第1天", signIn.fromUser.Uname)
			logic.PushToBulletSender(strMessage)
		}
	} else {
		sqlCommandBase = "insert into '%d' (uid,lastday,numberOfConsecutiveCheckInDays)  values (%d, %d, %d);"
		sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, signIn.fromUser.Uid, now.Unix(), 1)
		if _, err = db.Exec(sqlCommand); err != nil {
			db.Close()
			logx.Errorf("SQL执行失败：%s", err)
			return
		}
		strMessage := fmt.Sprintf("%s,连续签到第1天", signIn.fromUser.Uname)
		logic.PushToBulletSender(strMessage)
	}

	db.Close()
}

func (signIn *SignIn) Create() DanmuProcess {
	return new(SignIn)
}

func (signIn *SignIn) SetConfig(svcCtx *svc.ServiceContext) {
	signIn.svcCtx = svcCtx
}

func (signIn *SignIn) SetDanmu(content *string, user *message.User) {
	signIn.danmuContent = content
	signIn.fromUser = user
}
