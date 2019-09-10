package serial

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type ProtocolType int

const (
	ProtocolTypeJSON ProtocolType = iota + 1
	ProtocolTypeBINARY
)

type Marshaler interface {
	MarshalSerial() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalSerial([]byte) error
}

func Marshal(p ProtocolType, v interface{}) ([]byte, error) {
	var marshaler Marshaler

	switch p {
	case ProtocolTypeJSON:
		marshaler = &jsonImpl{v: v}
	case ProtocolTypeBINARY:
		marshaler = &binaryImpl{v: v}
	default:
		panic(fmt.Sprintf("unknown ProtocolType: %d", p))
	}

	return marshaler.MarshalSerial()
}

func Unmarshal(data []byte, v interface{}) error {
	var unmarshaler Unmarshaler

	p := GuessProtocolType(data)
	switch p {
	case ProtocolTypeJSON:
		unmarshaler = &jsonImpl{v: v}
	case ProtocolTypeBINARY:
		unmarshaler = &binaryImpl{v: v}
	}

	data, err := Unpack(p, data)
	if err != nil {
		return err
	}

	return unmarshaler.UnmarshalSerial(data)
}

// GuessProtocolType Guess raw data starts with JSON or BINARY protocol
// JSON is default
func GuessProtocolType(b []byte) ProtocolType {
	// starts with {0x0f, 0x02}
	if len(b) > 2 && bytes.Equal(b[:2], []byte{0x0f, 2}) {
		return ProtocolTypeBINARY
	}
	return ProtocolTypeJSON
}

// Pack Packs with protocol head, tail, data length
func Pack(p ProtocolType, b []byte) ([]byte, error) {
	out := new(bytes.Buffer)
	// magic code head
	out.WriteByte(0x0f)
	// protocol type
	out.WriteByte(byte(p))
	// data length
	if err := binary.Write(out, binary.BigEndian, uint16(len(b))); err != nil {
		return nil, err
	}
	// data
	out.Write(b)
	// valid (not implemented)
	out.WriteByte(0x00)
	// magic code tail
	out.WriteByte(0xf0)
	return out.Bytes(), nil
}

// Unpack Unpacks protocol head, tail, data length, get data back
func Unpack(p ProtocolType, b []byte) ([]byte, error) {
	if len(b) < 6 {
		return nil, errors.New("invalid data")
	}

	// find magic code head and valid
	idx := bytes.LastIndex(b, []byte{0x0f, byte(p)})
	if idx == -1 {
		return nil, errors.New("magic code head is not valid")
	}

	// find magic code tail and valid
	if !bytes.Equal(b[len(b)-2:], []byte{0x00, 0xf0}) {
		return nil, errors.New("magic code tail is not valid")
	}

	// got data
	data := b[idx:]

	// valid data length
	var shouldLen uint16
	if err := binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &shouldLen); err != nil {
		return nil, err
	}
	if len(data) != int(shouldLen+4+2) {
		return nil, errors.New("data length is not valid")
	}

	return data[4 : shouldLen+4], nil
}
