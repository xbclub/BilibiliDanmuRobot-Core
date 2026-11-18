package entity

import (
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
)

type FansMedalInfo struct {
	TargetID     int64  `json:"target_id"`
	MedalLevel   int64  `json:"medal_level"`
	MedalName    string `json:"medal_name"`
	IsLighted    int64  `json:"is_lighted"`
	GuardLevel   int64  `json:"guard_level"`
	AnchorRoomid int64  `json:"anchor_roomid"`
}

type InteractWordV2 struct {
	UID       int64
	Uname     string
	MsgType   int32
	FansMedal *FansMedalInfo
}

func ParseInteractWordV2(pbB64 string) (*InteractWordV2, error) {
	raw, err := decodeInteractBase64(pbB64)
	if err != nil {
		return nil, err
	}
	out := &InteractWordV2{}
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
			out.MsgType = int32(v)
			raw = raw[n:]
		case 9:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			medal, err := parseFansMedal(v)
			if err != nil {
				return nil, err
			}
			out.FansMedal = medal
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

func parseFansMedal(raw []byte) (*FansMedalInfo, error) {
	out := &FansMedalInfo{}
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		raw = raw[n:]
		switch num {
		case 1, 2, 8, 9, 12:
			v, n := consumeVarintField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			switch num {
			case 1:
				out.TargetID = int64(v)
			case 2:
				out.MedalLevel = int64(v)
			case 8:
				out.IsLighted = int64(v)
			case 9:
				out.GuardLevel = int64(v)
			case 12:
				out.AnchorRoomid = int64(v)
			}
			raw = raw[n:]
		case 3:
			v, n := consumeBytesField(raw, typ)
			if n < 0 {
				return nil, protowire.ParseError(n)
			}
			out.MedalName = string(v)
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

func consumeVarintField(raw []byte, typ protowire.Type) (uint64, int) {
	if typ != protowire.VarintType {
		return 0, -1
	}
	return protowire.ConsumeVarint(raw)
}

func consumeBytesField(raw []byte, typ protowire.Type) ([]byte, int) {
	if typ != protowire.BytesType {
		return nil, -1
	}
	return protowire.ConsumeBytes(raw)
}

func decodeInteractBase64(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		// 空白 payload 直接报错，避免下游收到 uid=0 的零值事件
		return nil, fmt.Errorf("decode pb: empty payload")
	}
	s = strings.TrimRight(s, "=")
	if raw, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return raw, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decode pb: %w", err)
	}
	return raw, nil
}
