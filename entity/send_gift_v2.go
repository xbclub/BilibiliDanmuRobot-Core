package entity

import (
	"math"

	"google.golang.org/protobuf/encoding/protowire"
)

type SendGiftV2Text struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Dmscore int    `json:"dmscore"`
		Pb      string `json:"pb"`
		Data    struct {
			Pb string `json:"pb"`
		} `json:"data"`
	} `json:"data"`
}

type SendGiftV2 struct {
	UID        int64
	Uname      string
	GuardLevel int64
	Gifts      []SendGiftV2Gift
	BlindGift  *SendGiftV2BlindGift
}

type SendGiftV2Gift struct {
	GiftID            int64
	GiftName          string
	Num               int
	GiftType          int64
	Price             int64
	DiscountPrice     int64
	TotalCoin         int64
	CoinType          string
	Tid               string
	Timestamp         int64
	SuperBatchGiftNum int64
	BatchComboID      string
	ComboResourcesID  int64
	ComboTotalCoin    int64
	ComboStayTime     int64
	Magnification     float32
}

type SendGiftV2BlindGift struct {
	OriginalGiftID   int64
	OriginalGiftName string
	BlindPrice       int64
}

// ParseSendGiftV2 解析 SEND_GIFT_V2 的 data.pb（base64 protobuf）。
// 字段号以实测抓包为准：外层 1=uid 2=uname 5=guard_level 9=blind_gift 10=gift_list 15=sender_uinfo；
// GiftItem 1=gift_id 2=gift_name 3=num 5=price 6=discount_price 7=total_coin 8=coin_type
// 9=tid 10=timestamp 11=super_batch_gift_num 12=batch_combo_id 13=combo_resources_id
// 14=combo_total_coin 15=combo_stay_time 16=magnification。
func ParseSendGiftV2(pbB64 string) (*SendGiftV2, error) {
	raw, err := decodeInteractBase64(pbB64)
	if err != nil {
		return nil, err
	}
	out := &SendGiftV2{}
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		raw = raw[n:]
		switch num {
		case 1:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			out.UID = int64(v)
			raw = raw[n:]
		case 2:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			out.Uname = string(v)
			raw = raw[n:]
		case 5:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			out.GuardLevel = int64(v)
			raw = raw[n:]
		case 9:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			blind, err := parseSendGiftV2BlindGift(v)
			if err != nil {
				return nil, err
			}
			out.BlindGift = blind
			raw = raw[n:]
		case 10:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			item, err := parseSendGiftV2Gift(v)
			if err != nil {
				return nil, err
			}
			out.Gifts = append(out.Gifts, item)
			raw = raw[n:]
		case 15:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			uid, uname := parseSendGiftV2Sender(v)
			if out.UID == 0 {
				out.UID = uid
			}
			if out.Uname == "" {
				out.Uname = uname
			}
			raw = raw[n:]
		default:
			n := protowire.ConsumeFieldValue(num, typ, raw)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			raw = raw[n:]
		}
	}
	return out, nil
}

func parseSendGiftV2Gift(raw []byte) (SendGiftV2Gift, error) {
	var out SendGiftV2Gift
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return out, protowire.ParseError(n)
		}
		raw = raw[n:]
		switch num {
		case 1:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.GiftID = int64(v)
			raw = raw[n:]
		case 2:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.GiftName = string(v)
			raw = raw[n:]
		case 3:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.Num = int(v)
			raw = raw[n:]
		case 4:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.GiftType = int64(v)
			raw = raw[n:]
		case 5:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.Price = int64(v)
			raw = raw[n:]
		case 6:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.DiscountPrice = int64(v)
			raw = raw[n:]
		case 7:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.TotalCoin = int64(v)
			raw = raw[n:]
		case 8:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.CoinType = string(v)
			raw = raw[n:]
		case 9:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.Tid = string(v)
			raw = raw[n:]
		case 10:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.Timestamp = int64(v)
			raw = raw[n:]
		case 11:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.SuperBatchGiftNum = int64(v)
			raw = raw[n:]
		case 12:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.BatchComboID = string(v)
			raw = raw[n:]
		case 13:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.ComboResourcesID = int64(v)
			raw = raw[n:]
		case 14:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.ComboTotalCoin = int64(v)
			raw = raw[n:]
		case 15:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.ComboStayTime = int64(v)
			raw = raw[n:]
		case 16:
			if typ != protowire.Fixed32Type {
				return out, protowire.ParseError(-1)
			}
			v, n := protowire.ConsumeFixed32(raw)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			out.Magnification = math.Float32frombits(v)
			raw = raw[n:]
		default:
			n := protowire.ConsumeFieldValue(num, typ, raw)
			if n < 0 {
				return out, protowire.ParseError(n)
			}
			raw = raw[n:]
		}
	}
	return out, nil
}

func parseSendGiftV2BlindGift(raw []byte) (*SendGiftV2BlindGift, error) {
	out := &SendGiftV2BlindGift{}
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		raw = raw[n:]
		switch num {
		case 2:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			out.OriginalGiftID = int64(v)
			raw = raw[n:]
		case 3:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			out.OriginalGiftName = string(v)
			raw = raw[n:]
		case 6:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			out.BlindPrice = int64(v)
			raw = raw[n:]
		default:
			n := protowire.ConsumeFieldValue(num, typ, raw)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			raw = raw[n:]
		}
	}
	return out, nil
}

func parseSendGiftV2Sender(raw []byte) (int64, string) {
	var uid int64
	var uname string
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return uid, uname
		}
		raw = raw[n:]
		switch num {
		case 1:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return uid, uname
			}
			uid = int64(v)
			raw = raw[n:]
		case 2:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return uid, uname
			}
			base := v
			for len(base) > 0 {
				bn, bt, bn2 := protowire.ConsumeTag(base)
				if bn2 < 0 {
					return uid, uname
				}
				base = base[bn2:]
				if bn == 1 && bt == protowire.BytesType {
					name, n2 := protowire.ConsumeBytes(base)
					if n2 < 0 {
						return uid, uname
					}
					uname = string(name)
					base = base[n2:]
				} else {
					n2 := protowire.ConsumeFieldValue(bn, bt, base)
					if n2 < 0 {
						return uid, uname
					}
					base = base[n2:]
				}
			}
		default:
			n := protowire.ConsumeFieldValue(num, typ, raw)
			if n < 0 {
				return uid, uname
			}
			raw = raw[n:]
		}
	}
	return uid, uname
}

func (g SendGiftV2Gift) numAndTotal() (int, int64) {
	layouts := [][2]int64{
		{int64(g.Num), g.TotalCoin},
		{g.SuperBatchGiftNum, g.DiscountPrice},
	}
	for _, layout := range layouts {
		num, total := layout[0], layout[1]
		if num > 0 && total > 0 && g.Price > 0 && total == g.Price*num {
			return int(num), total
		}
	}
	num := g.Num
	if num <= 0 {
		num = int(g.SuperBatchGiftNum)
	}
	if num <= 0 {
		num = 1
	}
	total := g.TotalCoin
	if total == 0 {
		total = g.DiscountPrice
	}
	if total == 0 {
		total = g.Price * int64(num)
	}
	return num, total
}

func (g *SendGiftV2) ToSendGiftText(item SendGiftV2Gift) *SendGiftText {
	num, total := item.numAndTotal()
	out := &SendGiftText{Cmd: "SEND_GIFT"}
	out.Data.UID = int(g.UID)
	out.Data.Uname = g.Uname
	out.Data.GuardLevel = int(g.GuardLevel)
	out.Data.Giftid = int(item.GiftID)
	out.Data.GiftName = item.GiftName
	out.Data.Gifttype = int(item.GiftType)
	out.Data.Num = num
	out.Data.Price = int(item.Price)
	out.Data.TotalCoin = int(total)
	out.Data.CoinType = item.CoinType
	out.Data.Tid = item.Tid
	out.Data.Timestamp = int(item.Timestamp)
	out.Data.SuperBatchGiftNum = int(item.SuperBatchGiftNum)
	out.Data.BatchComboID = item.BatchComboID
	out.Data.ComboResourcesID = int(item.ComboResourcesID)
	out.Data.ComboTotalCoin = int(item.ComboTotalCoin)
	out.Data.ComboStayTime = int(item.ComboStayTime)
	out.Data.Magnification = float64(item.Magnification)
	out.Data.Action = "投喂"
	if g.BlindGift != nil {
		out.Data.BlindGift.OriginalGiftId = int(g.BlindGift.OriginalGiftID)
		out.Data.BlindGift.OriginalGiftName = g.BlindGift.OriginalGiftName
		out.Data.BlindGift.OriginalGiftPrice = int(g.BlindGift.BlindPrice)
		if out.Data.GiftName == "" {
			out.Data.GiftName = g.BlindGift.OriginalGiftName
		}
	}
	return out
}
