package entity

import (
	"encoding/base64"
	"testing"

	"google.golang.org/protobuf/encoding/protowire"
)

func appendVarintField(b []byte, field int, v uint64) []byte {
	b = protowire.AppendTag(b, protowire.Number(field), protowire.VarintType)
	return protowire.AppendVarint(b, v)
}

func appendStringField(b []byte, field int, s string) []byte {
	b = protowire.AppendTag(b, protowire.Number(field), protowire.BytesType)
	return protowire.AppendString(b, s)
}

func appendBytesField(b []byte, field int, v []byte) []byte {
	b = protowire.AppendTag(b, protowire.Number(field), protowire.BytesType)
	return protowire.AppendBytes(b, v)
}

func TestParseSendGiftV2(t *testing.T) {
	var item []byte
	item = appendVarintField(item, 1, 35525)
	item = appendStringField(item, 2, "甜滋滋")
	item = appendVarintField(item, 3, 5)
	item = appendVarintField(item, 4, 2)
	item = appendVarintField(item, 5, 5000)
	item = appendVarintField(item, 6, 5000)
	item = appendVarintField(item, 7, 25000)
	item = appendStringField(item, 8, "gold")
	item = appendStringField(item, 9, "4815475693083221504")
	item = appendVarintField(item, 10, 1789043303)
	item = appendVarintField(item, 11, 1)
	item = appendStringField(item, 12, "batch:gift:combo_id:1")
	item = appendVarintField(item, 13, 10)
	item = appendVarintField(item, 14, 25000)
	item = appendVarintField(item, 15, 6)

	var top []byte
	top = appendVarintField(top, 1, 272404088)
	top = appendStringField(top, 2, "测试")
	top = appendBytesField(top, 10, item)

	got, err := ParseSendGiftV2(base64.StdEncoding.EncodeToString(top))
	if err != nil {
		t.Fatalf("ParseSendGiftV2() error = %v", err)
	}
	if got.UID != 272404088 || got.Uname != "测试" || len(got.Gifts) != 1 {
		t.Fatalf("unexpected sender/gifts: %+v", got)
	}
	send := got.ToSendGiftText(got.Gifts[0])
	if send.Data.GiftName != "甜滋滋" || send.Data.Num != 5 || send.Data.Price != 5000 || send.Data.TotalCoin != 25000 {
		t.Fatalf("unexpected converted gift: %+v", send.Data)
	}
	if send.Data.UID != 272404088 || send.Data.Uname != "测试" || send.Data.Giftid != 35525 {
		t.Fatalf("unexpected converted sender: %+v", send.Data)
	}
}

func TestParseSendGiftV2AlternateLayout(t *testing.T) {
	var item []byte
	item = appendVarintField(item, 1, 33988)
	item = appendStringField(item, 2, "人气票")
	item = appendVarintField(item, 5, 100)
	item = appendVarintField(item, 6, 300)
	item = appendStringField(item, 8, "gold")
	item = appendStringField(item, 9, "4798847560141970944")
	item = appendVarintField(item, 10, 1785078847)
	item = appendVarintField(item, 11, 3)
	item = appendStringField(item, 12, "batch:gift:combo_id:2")
	item = appendVarintField(item, 13, 10)

	var top []byte
	top = appendVarintField(top, 1, 26928797)
	top = appendStringField(top, 2, "mikufilck")
	top = appendBytesField(top, 10, item)

	got, err := ParseSendGiftV2(base64.StdEncoding.EncodeToString(top))
	if err != nil {
		t.Fatalf("ParseSendGiftV2() error = %v", err)
	}
	send := got.ToSendGiftText(got.Gifts[0])
	if send.Data.Num != 3 || send.Data.TotalCoin != 300 || send.Data.Price != 100 || send.Data.GiftName != "人气票" {
		t.Fatalf("alternate layout not handled: %+v", send.Data)
	}
}

func TestParseSendGiftV2BlindGift(t *testing.T) {
	var blind []byte
	blind = appendVarintField(blind, 2, 13000)
	blind = appendStringField(blind, 3, "粉丝团灯牌")
	blind = appendVarintField(blind, 6, 500)

	var item []byte
	item = appendVarintField(item, 1, 13000)
	item = appendStringField(item, 2, "盲盒")
	item = appendVarintField(item, 3, 1)
	item = appendVarintField(item, 5, 500)
	item = appendVarintField(item, 7, 500)

	var top []byte
	top = appendVarintField(top, 1, 1234)
	top = appendStringField(top, 2, "盲盒用户")
	top = appendBytesField(top, 9, blind)
	top = appendBytesField(top, 10, item)

	got, err := ParseSendGiftV2(base64.StdEncoding.EncodeToString(top))
	if err != nil {
		t.Fatalf("ParseSendGiftV2() error = %v", err)
	}
	send := got.ToSendGiftText(got.Gifts[0])
	if send.Data.BlindGift.OriginalGiftName != "粉丝团灯牌" || send.Data.BlindGift.OriginalGiftPrice != 500 || send.Data.BlindGift.OriginalGiftId != 13000 {
		t.Fatalf("blind gift not handled: %+v", send.Data.BlindGift)
	}
}

func TestParseSendGiftV2Invalid(t *testing.T) {
	if _, err := ParseSendGiftV2("!!!"); err == nil {
		t.Fatal("expected an error for invalid base64")
	}
}
