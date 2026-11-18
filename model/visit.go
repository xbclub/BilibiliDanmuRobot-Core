package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VisitModel interface {
	Touch(ctx context.Context, uid, now int64) (first bool, err error)
}

type defaultVisitModel struct {
	conn  *gorm.DB
	table string
}

type VisitBase struct {
	Uid       int64 `gorm:"primaryKey;autoIncrement:false"`
	FirstSeen int64
	LastSeen  int64
	VisitCnt  int64
}

// visitRetention 进场记录保留窗口，建表/重连时清理长期未访问的行
const visitRetention = 90 * 24 * 3600

func NewVisitModel(conn *gorm.DB, roomID int64) VisitModel {
	table := fmt.Sprintf("visit_%v", roomID)
	if err := conn.Table(table).AutoMigrate(&VisitBase{}); err != nil {
		// 建表失败说明数据库已不可用，静默继续只会让三档欢迎一直走回退路径
		logx.Must(err)
	}
	// 清理保留窗口外的旧记录，避免 per-room 表无限增长
	if err := conn.Table(table).
		Where("last_seen < ?", time.Now().Unix()-visitRetention).
		Delete(&VisitBase{}).Error; err != nil {
		logx.Errorf("清理进场记录失败 table=%s: %v", table, err)
	}
	return &defaultVisitModel{
		conn:  conn,
		table: table,
	}
}

func (m *defaultVisitModel) Touch(ctx context.Context, uid, now int64) (bool, error) {
	var first bool
	err := m.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).
			Table(m.table).
			Create(&VisitBase{Uid: uid, FirstSeen: now, LastSeen: now, VisitCnt: 1})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 1 {
			first = true
			return nil
		}
		return tx.Table(m.table).Where("uid = ?", uid).Updates(map[string]interface{}{
			"last_seen": now,
			"visit_cnt": gorm.Expr("visit_cnt + 1"),
		}).Error
	})
	return first, err
}
