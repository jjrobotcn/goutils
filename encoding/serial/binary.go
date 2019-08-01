package serial

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

const tagName = "binary"

type binaryImpl struct {
	v   interface{}
	buf *bytes.Buffer
}

func (b *binaryImpl) UnmarshalSerial(p []byte) error {
	t := reflect.TypeOf(b.v)
	if t.Kind() != reflect.Ptr {
		return errors.New("require ptr")
	}
	t = t.Elem()

	// p position
	var pos int
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		fieldKind := tField.Type.Kind()
		fieldTag := tField.Tag.Get(tagName)

		// ignore tag
		if fieldTag == "-" {
			continue
		}

		v := reflect.ValueOf(b.v).Elem()
		vField := v.FieldByName(tField.Name)

		// ignore private field
		if !vField.CanInterface() {
			continue
		}

		if fieldTag == "" {
			switch fieldKind {
			case reflect.Uint8:
				vField.Set(reflect.ValueOf(p[pos]))
				pos++
			case reflect.Int8:
				vField.Set(reflect.ValueOf(int8(p[pos])))
				pos++
			case reflect.Uint16:
				bb := p[pos : pos+2]
				u16 := binary.BigEndian.Uint16(bb)
				vField.Set(reflect.ValueOf(u16))
				pos += 2
			case reflect.Int16:
				bb := p[pos : pos+2]
				u16 := binary.BigEndian.Uint16(bb)
				vField.Set(reflect.ValueOf(int16(u16)))
				pos += 2
			case reflect.Bool:
				buf := new(bytes.Buffer)
				buf.WriteByte(p[pos])
				var bb bool
				if err := binary.Read(buf, binary.BigEndian, &bb); err != nil {
					return err
				}
				vField.SetBool(bb)
			}
		}

	}

	return nil
}

func (b *binaryImpl) MarshalSerial() ([]byte, error) {
	// init buffer
	b.buf = new(bytes.Buffer)

	t := reflect.TypeOf(b.v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		fieldKind := tField.Type.Kind()
		fieldTag := tField.Tag.Get(tagName)

		// ignore tag
		if fieldTag == "-" {
			continue
		}

		v := reflect.ValueOf(b.v)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		vField := v.FieldByName(tField.Name)

		// ignore private
		if !vField.CanInterface() {
			continue
		}

		// default
		if fieldTag == "" {
			if err := b.marshalDefault(fieldKind, vField); err != nil {
				return nil, err
			}
		}
	}

	return pack(ProtocolTypeBINARY, b.buf.Bytes())
}

func (b *binaryImpl) marshalDefault(kind reflect.Kind, v reflect.Value) error {
	switch kind {
	case reflect.Uint8:
		return binary.Write(b.buf, binary.BigEndian, uint8(v.Uint()))
	case reflect.Int8:
		return binary.Write(b.buf, binary.BigEndian, int8(v.Int()))
	case reflect.Uint16:
		return binary.Write(b.buf, binary.BigEndian, uint16(v.Uint()))
	case reflect.Int16:
		return binary.Write(b.buf, binary.BigEndian, int16(v.Int()))
	case reflect.Bool:
		return binary.Write(b.buf, binary.BigEndian, v.Bool())
	}
	return fmt.Errorf("kind %s is not support", kind)
}
