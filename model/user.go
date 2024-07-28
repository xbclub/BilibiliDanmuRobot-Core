package model

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type (
	SignInModel interface {
		Insert(ctx context.Context, tx *gorm.DB, data *SingInBase) error
		FindOne(ctx context.Context, id int64) (*SingInBase, error)
		UpdateCount(ctx context.Context, uid int64) error
	}
	defaultSingInModel struct {
		conn  *gorm.DB
		table string
	}
	SingInBase struct {
		ID      int64 `gorm:"primaryKey;autoIncrement"`
		Uid     int64
		LastDay int64
		Count   int64
	}
)

func NewSignInModel(conn *gorm.DB, RoomID int64) SignInModel {
	err := conn.Table(fmt.Sprintf("room_%v", RoomID)).AutoMigrate(&SingInBase{})
	if err != nil {
		logx.Error(err)
	}
	return &defaultSingInModel{
		conn:  conn,
		table: fmt.Sprintf("room_%v", RoomID),
	}
}

func (m *defaultSingInModel) Insert(ctx context.Context, tx *gorm.DB, data *SingInBase) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Table(m.table).Save(&data).Error
	return err
}
func (m *defaultSingInModel) FindOne(ctx context.Context, uid int64) (*SingInBase, error) {
	var resp SingInBase
	err := m.conn.WithContext(ctx).Table(m.table).Model(&SingInBase{}).Where("uid = ?", uid).Take(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultSingInModel) UpdateCount(ctx context.Context, uid int64) error {
	err := m.conn.WithContext(ctx).Table(m.table).Model(&SingInBase{}).Where("uid = ?", uid).UpdateColumn("count", gorm.Expr("count + ?", 1)).UpdateColumn("last_day", carbon.Now(carbon.Local).Timestamp()).Error
	return err
}
