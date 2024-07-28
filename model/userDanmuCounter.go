package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

type (
	DanmuCntModel interface {
		Insert(ctx context.Context, tx *gorm.DB, data *DanmuCntBase) error
		FindOne(ctx context.Context, id int64, date string) (*DanmuCntBase, error)
		UpdateCount(ctx context.Context, uid int64) error
		GetRecent3DayRecords(ctx context.Context, uid int64) ([]DanmuCntBase, error)
		GetDateStr(daysFromToday int) string
	}
	defaultDanmuCntModel struct {
		conn  *gorm.DB
		table string
	}
	DanmuCntBase struct {
		ID    int64 `gorm:"primaryKey;autoIncrement"`
		Uid   int64
		Date  string
		Count int64
	}
)

func NewDanmuCntModel(conn *gorm.DB, RoomID int64) DanmuCntModel {
	err := conn.Table(fmt.Sprintf("room_%v", RoomID)).AutoMigrate(&DanmuCntBase{})
	if err != nil {
		logx.Error(err)
	}
	return &defaultDanmuCntModel{
		conn:  conn,
		table: fmt.Sprintf("room_%v", RoomID),
	}
}

func (m *defaultDanmuCntModel) GetDateStr(daysFromToday int) string {
	beijingLocation, err1 := time.LoadLocation("Asia/Shanghai")
	if err1 != nil {
		logx.Error(err1)
	}
	dateStr := time.Now().In(beijingLocation).AddDate(0, 0, -daysFromToday).Format("2006-01-02")
	logx.Info(dateStr)
	return dateStr
}

func (m *defaultDanmuCntModel) Insert(ctx context.Context, tx *gorm.DB, data *DanmuCntBase) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Table(m.table).Save(&data).Error
	return err
}

func (m *defaultDanmuCntModel) FindOne(ctx context.Context, uid int64, date string) (*DanmuCntBase, error) {
	var resp DanmuCntBase
	err := m.conn.WithContext(ctx).Table(m.table).Model(&DanmuCntBase{}).Where("uid = ? AND date = ?", uid, date).Take(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultDanmuCntModel) GetRecent3DayRecords(ctx context.Context, uid int64) ([]DanmuCntBase, error) {
	var resp []DanmuCntBase
	endDate := m.GetDateStr(0)
	startDate := m.GetDateStr(2)
	err := m.conn.WithContext(ctx).Table(m.table).Model(&DanmuCntBase{}).Where("uid = ? AND date BETWEEN ? AND ?", uid, startDate, endDate).Order("date asc").Take(&resp).Error
	switch err {
	case nil:
		return resp, nil
	case ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultDanmuCntModel) UpdateCount(ctx context.Context, uid int64) error {
	today := m.GetDateStr(0)
	err := m.conn.WithContext(ctx).Table(m.table).Model(&SingInBase{}).Where("uid = ? AND date = ?", uid, today).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error
	return err
}
