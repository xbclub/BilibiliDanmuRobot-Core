package model

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestVisitModelTouch(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "visit.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	defer sqlDB.Close()
	m := NewVisitModel(db, 12345)

	first, err := m.Touch(context.Background(), 1001, 1000)
	if err != nil {
		t.Fatalf("first Touch: %v", err)
	}
	if !first {
		t.Fatal("first Touch should report first visit")
	}

	first, err = m.Touch(context.Background(), 1001, 2000)
	if err != nil {
		t.Fatalf("second Touch: %v", err)
	}
	if first {
		t.Fatal("second Touch should not report first visit")
	}

	var row VisitBase
	if err := db.Table("visit_12345").Where("uid = ?", 1001).Take(&row).Error; err != nil {
		t.Fatalf("read visit row: %v", err)
	}
	if row.FirstSeen != 1000 || row.LastSeen != 2000 || row.VisitCnt != 2 {
		t.Fatalf("unexpected row: %+v", row)
	}

	first, err = m.Touch(context.Background(), 1002, 3000)
	if err != nil {
		t.Fatalf("other user Touch: %v", err)
	}
	if !first {
		t.Fatal("a different uid should be a first visit")
	}
}
