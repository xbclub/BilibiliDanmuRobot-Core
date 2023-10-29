package svc

import (
	"database/sql"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type ServiceContext struct {
	Config       config.Config
	OtherSideUid map[int64]bool
	TopUid       map[int64]bool
	Db           struct {
		TableName string
		Db        *sql.DB
	}
}

func NewServiceContext(c config.Config) *ServiceContext {
	tablename, db, err := connectSqlite(c)
	if err != nil {
		logx.Error(err)
	}
	return &ServiceContext{
		OtherSideUid: make(map[int64]bool),
		TopUid:       make(map[int64]bool),
		Db: struct {
			TableName string
			Db        *sql.DB
		}{TableName: tablename, Db: db},
		Config: c,
	}
}

func connectSqlite(c config.Config) (tablename string, db *sql.DB, err error) {
	strRoomID := strconv.Itoa(c.RoomId) // s is "42"
	tableName := "signin_" + strRoomID
	dbFile := fmt.Sprintf("%s/%s?_pragma=busy_timeout(5000)", c.DBPath, c.DBName)

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		logx.Errorf("SQL执行失败：%s", err)
		return tableName, db, err
	}

	sqlCommandBase := `
		create table if not exists '%s'(uid,lastday,numberOfConsecutiveCheckInDays)`

	sqlCommand := fmt.Sprintf(sqlCommandBase, tableName)
	if _, err := db.Exec(sqlCommand); err != nil {
		logx.Errorf("SQL执行失败：%s", err)
		return tableName, db, err
	}
	return tableName, db, nil
}
