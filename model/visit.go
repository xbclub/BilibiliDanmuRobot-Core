package model

import (
	"context"
	"fmt"

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

func NewVisitModel(conn *gorm.DB, roomID int64) VisitModel {
	table := fmt.Sprintf("visit_%v", roomID)
	if err := conn.Table(table).AutoMigrate(&VisitBase{}); err != nil {
		logx.Error(err)
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
