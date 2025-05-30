package packet

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/andybalholm/brotli"
	log "github.com/zeromicro/go-zero/core/logx"
)

const (
	Plain = iota
	Popularity
	Zlib
	Brotli
)

const (
	HandShake = iota
	HandShakeResponse
	HeartBeat
	HeartBeatResponse
	_
	Notification
	_
	RoomEnter
	RoomEnterResponse
)

type Packet struct {
	PacketLength    int // PacketLength 在 build 时会计算
	HeaderLength    int // HeaderLength 大概是固定值 16
	ProtocolVersion uint16
	Operation       uint32
	SequenceID      int
	Body            []byte
}

func NewPacket(protocolVersion uint16, operation uint32, body []byte) Packet {
	return Packet{
		ProtocolVersion: protocolVersion,
		Operation:       operation,
		Body:            body,
	}
}

// NewPlainPacket 构造新的 Plain 包
// 对外暴露的方法中 operation 全部使用int
func NewPlainPacket(operation int, body []byte) Packet {
	return NewPacket(1, uint32(operation), body)
}

func NewPacketFromBytes(data []byte) Packet {
	if len(data) == 0 {
		log.Error("error packet")
		return Packet{ProtocolVersion: 1, Operation: uint32(3), Body: make([]byte, 0)}
	}

	packLen := binary.BigEndian.Uint32(data[0:4])

	if int(packLen) == 0 || int(packLen) != len(data) {
		log.Error("error packet")
		return Packet{ProtocolVersion: 1, Operation: uint32(3), Body: make([]byte, 0)}
	}

	pv := binary.BigEndian.Uint16(data[6:8])
	op := binary.BigEndian.Uint32(data[8:12])
	body := data[16:packLen]
	packet := NewPacket(pv, op, body)
	return packet
}

func (p Packet) Parse() []Packet {
	switch p.ProtocolVersion {
	case Popularity:
		fallthrough
	case Plain:
		return []Packet{p}
	case Zlib:
		z, err := zlibParser(p.Body)
		if err != nil {
			log.Error("zlib error", err)
		}
		return Slice(z)
	case Brotli:
		b, err := brotliParser(p.Body)
		if err != nil {
			log.Error("brotli error", err)
		}
		return Slice(b)
	default:
		log.Error("unknown protocolVersion")
	}
	return nil
}

func (p *Packet) Unmarshal(v interface{}) error {
	return json.Unmarshal(p.Body, v)
}

func (p *Packet) Build() []byte {
	// default packet with some constants
	// raw header length: buf[5] = 16
	// seq: buf[16] = 1
	rawBuf := []byte{0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// protocol version
	binary.BigEndian.PutUint16(rawBuf[6:], p.ProtocolVersion)

	// operation code
	binary.BigEndian.PutUint32(rawBuf[8:], p.Operation)

	// append payload
	rawBuf = append(rawBuf, p.Body...)

	// payload length
	binary.BigEndian.PutUint32(rawBuf, uint32(len(rawBuf)))
	return rawBuf
}

// DecodePacket Decode
func DecodePacket(data []byte) Packet {
	return NewPacketFromBytes(data)
}

// EncodePacket Encode
func EncodePacket(packet Packet) []byte {
	return packet.Build()
}

func Slice(data []byte) []Packet {
	var packets []Packet
	total := len(data)
	cursor := 0
	for cursor < total {
		packLen := int(binary.BigEndian.Uint32(data[cursor : cursor+4]))
		packets = append(packets, DecodePacket(data[cursor:cursor+packLen]))
		cursor += packLen
	}
	return packets
}

func zlibParser(b []byte) ([]byte, error) {
	var rdBuf []byte
	zr, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	rdBuf, err = io.ReadAll(zr)
	return rdBuf, nil
}

func brotliParser(b []byte) ([]byte, error) {
	zr := brotli.NewReader(bytes.NewReader(b))
	rdBuf, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return rdBuf, nil
}
