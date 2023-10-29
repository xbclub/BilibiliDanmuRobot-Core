package danmuProcess

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/Masterminds/squirrel"
	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-module/carbon/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SignIn struct {
	danmuContent *string
	fromUser     *message.User
	svcCtx       *svc.ServiceContext
	db           *sql.DB
	tableName    string
}

func (signIn *SignIn) DoDanmuProcess() {
	if *signIn.danmuContent != "签到" {
		return
	}

	//strUid := strconv.Itoa(signIn.fromUser.Uid)
	//sqlCommandBase = "select * from '%d' where uid=%d;"
	//sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, signIn.fromUser.Uid)
	users := squirrel.Select("*").From(signIn.tableName).Where("uid = ?", signIn.fromUser.Uid)

	sqlCommand, _, err := users.ToSql()
	if err != nil {
		logx.Errorf("SQL执行失败：%s", err)
		return
	}

	stmt, err := signIn.db.Prepare(sqlCommand)

	if err != nil {
		logx.Errorf("SQL执行失败：%s", err)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(signIn.fromUser.Uid)
	//if err != nil {
	//	logx.Errorf("SQL执行失败：%s", err)
	//	return
	//}

	var uid, lastday int
	var numberOfConsecutiveCheckInDays int64

	// 获取当前时间
	//now := time.Now().In(time.Local)
	now := carbon.Now(carbon.Shanghai)
	err = row.Scan(&uid, &lastday, &numberOfConsecutiveCheckInDays)
	if err == nil {
		//sqlCommandBase = "update '%d' set lastday=%d, numberOfConsecutiveCheckInDays=%d where uid=%d;"
		// 将时间戳转换为时间对象（中国时区）
		//lastTime := time.Unix(numberOfConsecutiveCheckInDays, 0).In(time.Local)
		lastTime := carbon.CreateFromTimestamp(numberOfConsecutiveCheckInDays, carbon.Shanghai)
		if now.Year() == lastTime.Year() && now.Month() == lastTime.Month() && now.Day() == lastTime.Day()+1 {
			//sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, now.Unix(), numberOfConsecutiveCheckInDays+1, signIn.fromUser.Uid)
			update := squirrel.Update("").Table(signIn.tableName).Set("lastday", now.Timestamp()).Set("numberOfConsecutiveCheckInDays", numberOfConsecutiveCheckInDays+1).Where(squirrel.Eq{"uid": signIn.fromUser.Uid})
			sqlCommand, _, err := update.ToSql()
			if err != nil {
				logx.Errorf("SQL执行失败：%s", err)
				return
			}

			stmt, err := signIn.db.Prepare(sqlCommand)
			if err != nil {
				logx.Errorf("SQL执行失败：%s", err)
				return
			}
			defer stmt.Close()

			if _, err = stmt.Exec(now.Timestamp(), numberOfConsecutiveCheckInDays+1, signIn.fromUser.Uid); err != nil {
				logx.Errorf("SQL执行失败：%s", err)
				return
			}

			strMessage := fmt.Sprintf("%s,连续签到第%d天", signIn.fromUser.Uname, numberOfConsecutiveCheckInDays+1)
			logic.PushToBulletSender(strMessage)
		} else if now.Year() == lastTime.Year() && now.Month() == lastTime.Month() && now.Day() != lastTime.Day() {
			//sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, now.Unix(), 1, signIn.fromUser.Uid)
			update := squirrel.Update("").Table(signIn.tableName).Set("lastday", now.Timestamp()).Set("numberOfConsecutiveCheckInDays", 1).Where(squirrel.Eq{"uid": signIn.fromUser.Uid})

			sqlCommand, _, err := update.ToSql()
			if err != nil {
				logx.Errorf("SQL执行失败：%s", err)
				return
			}

			stmt, err := signIn.db.Prepare(sqlCommand)
			if err != nil {
				logx.Errorf("SQL执行失败：%s", err)
				return
			}
			defer stmt.Close()

			if _, err = stmt.Exec(now.Timestamp(), 1, signIn.fromUser.Uid); err != nil {
				logx.Errorf("SQL执行失败：%s", err)
				return
			}

			strMessage := fmt.Sprintf("%s,连续签到第1天", signIn.fromUser.Uname)
			logic.PushToBulletSender(strMessage)
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		//sqlCommandBase = "insert into '%d' (uid,lastday,numberOfConsecutiveCheckInDays)  values (%d, %d, %d);"
		//sqlCommand = fmt.Sprintf(sqlCommandBase, signIn.svcCtx.Config.RoomId, signIn.fromUser.Uid, now.Unix(), 1)
		insert := squirrel.Insert(signIn.tableName).Columns("uid", "lastday", "numberOfConsecutiveCheckInDays").Values(signIn.fromUser.Uid, now.Timestamp(), 1)

		sqlCommand, _, err := insert.ToSql()
		if err != nil {
			logx.Errorf("SQL执行失败：%s", err)
			return
		}

		stmt, err := signIn.db.Prepare(sqlCommand)
		if err != nil {
			logx.Errorf("SQL执行失败：%s", err)
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(signIn.fromUser.Uid, now.Timestamp(), 1); err != nil {
			logx.Errorf("SQL执行失败：%s", err)
			return
		}
		strMessage := fmt.Sprintf("%s,连续签到第1天", signIn.fromUser.Uname)
		logic.PushToBulletSender(strMessage)
	} else {
		logx.Errorf("SQL执行失败：%s", err)
		return
	}
}

func (signIn *SignIn) Create() DanmuProcess {
	return new(SignIn)
}

func (signIn *SignIn) SetConfig(svcCtx *svc.ServiceContext) {
	signIn.svcCtx = svcCtx
	signIn.db = svcCtx.Db.Db
	signIn.tableName = svcCtx.Db.TableName
}

func (signIn *SignIn) SetDanmu(content *string, user *message.User) {
	signIn.danmuContent = content
	signIn.fromUser = user
}
