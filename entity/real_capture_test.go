package entity

import "testing"

// 真实抓包样本（来自 danmu-recoder 项目 recoder-logs/room1878681152/20261107_184748.jsonl，
// 2026-11-07 B 站线上事件），用于回归验证 pb 字段号映射，补齐合成样本的循环自证缺口：
// 两位进场用户均未佩戴粉丝牌（顶层无 field 9、uinfo 内也无 medal），FansMedal 应为 nil；
// 礼物事件的名称/数量/单价/总价/折扣价/coin 类型/combo id 相互自洽。

const realInteractPbNoMedal1 = "CKfL04oFEghIRW5jT3JMRSIBASgBMMDM6f8GON2DvNcGQN6V/aybNGIAeJTXpNbZsu/yGJoBALIBRwiny9OKBRI7CghIRW5jT3JMRRIvaHR0cHM6Ly9pMC5oZHNsYi5jb20vYmZzL2ZhY2UvbWVtYmVyL25vZmFjZS5qcGciADIAugEAwgEA"

const realInteractPbNoMedal2 = "CJPk27gBEgzmlq3lio3nhKHmiJEiAQEoATDAzOn/Bjjwg7zXBkCapP6smzRiAHjsm7qhnbPv8hiaAQCyAUsIk+TbuAESPwoM5pat5YqN54Sh5oiREi9odHRwczovL2kwLmhkc2xiLmNvbS9iZnMvZmFjZS9tZW1iZXIvbm9mYWNlLmpwZyIAMgC6AQDCAQA="

const realGiftPbYuanqi = "CMaVgLnVsqYGEhjmoYPoirHkvp3ml6fnrJHmmKXpo47jg5oaSmh0dHBzOi8vaTIuaGRzbGIuY29tL2Jmcy9mYWNlLzAyMTJlYWMyNzUzZDMzYzcwZTAzYjQyZjFmMDE1OWFjYmE3ZWFmYmUuanBnIgcjMDBEMUYxKANCKgimloDr7remBigpMgnoiJLmnI3ogos4i8L9B0DAgYMGSMCBgwZQwIGDBkogCLwBEMOaAhoM6Zi/54uX55uy55uSKgbniIblh7owqEZSgQUIxJoCEgzlhYPmsJTmu6HmoLwYASACKKhGMKhGOKhGQgRnb2xkShM0ODM2NDY4OTkyMTc0NjA0Mjg4UO+DvNcGWAFiJDMyZTVhMDZkLTcxMTktNDc4MC1hMDAzLTc4OTc1ZTk4MmIwOWgKcKhGeAWFAQAAgD+IAQGSAQbmipXlloLAAa6ezQHqAR0KFeiIkuiKmeiVvnNvdWZmbGVfcGx1cxDgms73BooC3gEI4JrO9wYS1QEKFeiIkuiKmeiVvnNvdWZmbGVfcGx1cxJKaHR0cHM6Ly9pMi5oZHNsYi5jb20vYmZzL2ZhY2UvN2Y2ZDc3NzI4OWMxZjk3ZjZjM2YzZjZhNzZmZjNmMGZiOTIzYmI0Mi5qcGcyYwoV6IiS6IqZ6JW+c291ZmZsZV9wbHVzEkpodHRwczovL2kyLmhkc2xiLmNvbS9iZnMvZmFjZS83ZjZkNzc3Mjg5YzFmOTdmNmMzZjNmNmE3NmZmM2YwZmI5MjNiYjQyLmpwZzoLIP///////////wGSAgCaAuoBCkpodHRwczovL3MxLmhkc2xiLmNvbS9iZnMvbGl2ZS8zNTYxZDMxNDJmZmEzNzY2ZjhiZGFkZGJmMjk2ZThjZGJmNjAyOTQ2LnBuZxJLaHR0cHM6Ly9pMC5oZHNsYi5jb20vYmZzL2xpdmUvNzQxZGQxZGQ3ZDIzNThiYmUyZjU2ZjJhNWY3NjA0NzhhZTEwZmU2YS53ZWJwGIcyIAEqSmh0dHBzOi8vaTAuaGRzbGIuY29tL2Jmcy9saXZlLzZiMDhiZDA2MzM1OTVkYzk0MzY2YTI5MGUyMjBkMGYxNWZmMzkwMmQuZ2lmoAKoRqoCAFgBagIIKXqcAwjGlYC51bKmBhLbAQoY5qGD6Iqx5L6d5pen56yR5pil6aOO44OaEkpodHRwczovL2kyLmhkc2xiLmNvbS9iZnMvZmFjZS8wMjEyZWFjMjc1M2QzM2M3MGUwM2I0MmYxZjAxNTlhY2JhN2VhZmJlLmpwZzJmChjmoYPoirHkvp3ml6fnrJHmmKXpo47jg5oSSmh0dHBzOi8vaTIuaGRzbGIuY29tL2Jmcy9mYWNlLzAyMTJlYWMyNzUzZDMzYzcwZTAzYjQyZjFmMDE1OWFjYmE3ZWFmYmUuanBnOgsg////////////ARqyAQoJ5Ly05riF5ae/ECUYi8L9ByCEof8HKP/RnwMwi8L9B0gBUOCazvcGWANg2bUKakpodHRwczovL2kwLmhkc2xiLmNvbS9iZnMvbGl2ZS80ODM2MGM4ZjNiN2RlODAzMWU4NmZmMWVmNGEyZGZjMGVjMmE2MWMyLnBuZ3oJIzRDN0RGRjk5ggEJIzRDN0RGRjk5igEHIzU4QTFGOJIBByNGRkZGRkaaAQkjNEM3REZGRTY="

func TestParseInteractWordV2RealCapture(t *testing.T) {
	ev, err := ParseInteractWordV2(realInteractPbNoMedal1)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if ev.UID != 1364518311 || ev.Uname != "HEncOrLE" || ev.MsgType != 1 {
		t.Fatalf("interact1 解码不符: uid=%d uname=%q msgType=%d", ev.UID, ev.Uname, ev.MsgType)
	}
	if ev.FansMedal != nil {
		t.Fatalf("未佩戴粉丝牌应为 nil, got %+v", ev.FansMedal)
	}

	ev, err = ParseInteractWordV2(realInteractPbNoMedal2)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if ev.UID != 387379731 || ev.Uname != "断劍無我" || ev.MsgType != 1 {
		t.Fatalf("interact2 解码不符: uid=%d uname=%q msgType=%d", ev.UID, ev.Uname, ev.MsgType)
	}
	if ev.FansMedal != nil {
		t.Fatalf("未佩戴粉丝牌应为 nil, got %+v", ev.FansMedal)
	}
}

func TestParseSendGiftV2RealCapture(t *testing.T) {
	g, err := ParseSendGiftV2(realGiftPbYuanqi)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if g.UID != 3546566411422406 || g.Uname != "桃花依旧笑春风ペ" || len(g.Gifts) != 1 {
		t.Fatalf("外层解码不符: uid=%d uname=%q gifts=%d", g.UID, g.Uname, len(g.Gifts))
	}
	item := g.Gifts[0]
	if item.GiftID != 36164 || item.GiftName != "元气满格" || item.Num != 1 {
		t.Fatalf("礼物解码不符: id=%d name=%q num=%d", item.GiftID, item.GiftName, item.Num)
	}
	if item.Price != 9000 || item.DiscountPrice != 9000 || item.TotalCoin != 9000 || item.CoinType != "gold" {
		t.Fatalf("金额解码不符: price=%d discount=%d total=%d coin=%q",
			item.Price, item.DiscountPrice, item.TotalCoin, item.CoinType)
	}
	if item.BatchComboID != "32e5a06d-7119-4780-a003-78975e982b09" {
		t.Fatalf("combo id 解码不符: %q", item.BatchComboID)
	}

	s := g.ToSendGiftText(item)
	if s.Data.UID != 3546566411422406 || s.Data.Uname != "桃花依旧笑春风ペ" ||
		s.Data.GiftName != "元气满格" || s.Data.Num != 1 ||
		s.Data.Price != 9000 || s.Data.TotalCoin != 9000 || s.Data.CoinType != "gold" {
		t.Fatalf("ToSendGiftText 转换不符: %+v", s.Data)
	}
}
