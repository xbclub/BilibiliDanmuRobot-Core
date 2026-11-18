package model

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VisitModel interface {
	Touch(ctx context.Context, uid, now int64) (first bool, err error)
}

type defaultVisitModel struct {
	conn        *gorm.DB
	table       string
	lastCleanNs atomic.Int64
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
		// 与兄弟模型策略一致：记录错误返回 nil，三档欢迎判 nil 自动回退旧逻辑，
		// 不在换房 reload 路径上杀死整个进程
		logx.Errorf("visit 表迁移失败 table=%s: %v", table, err)
		return nil
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
	m.maybeClean(now)
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

// maybeClean 每 24 小时清理一次保留窗口外的旧记录（CAS 保证并发只触发一次）
func (m *defaultVisitModel) maybeClean(nowUnix int64) {
	last := m.lastCleanNs.Load()
	if last != 0 && time.Since(time.Unix(0, last)) < 24*time.Hour {
		return
	}
	if !m.lastCleanNs.CompareAndSwap(last, time.Now().UnixNano()) {
		return
	}
	if err := m.conn.Table(m.table).
		Where("last_seen < ?", nowUnix-visitRetention).
		Delete(&VisitBase{}).Error; err != nil {
		logx.Errorf("周期清理进场记录失败 table=%s: %v", m.table, err)
	}
}
