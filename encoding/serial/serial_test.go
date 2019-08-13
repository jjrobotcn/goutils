package serial

import (
	"testing"
)

func TestMarshalUnmarshalJson(t *testing.T) {
	type TestJson struct {
		private bool
		Bool    bool    `json:"bool"`
		String  string  `json:"string"`
		Int     int     `json:"int"`
		Float32 float32 `json:"float32"`
	}

	j := &TestJson{
		private: true,
		Bool:    true,
		String:  "string",
		Int:     1,
		Float32: 3.1415,
	}
	b, err := Marshal(ProtocolTypeJSON, j)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var check TestJson
	err = Unmarshal(b, &check)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if check.private != false || check.Bool != j.Bool || check.String != j.String || check.Int != j.Int || check.Float32 != j.Float32 {
		t.Errorf("want: %#v, got: %#v", j, check)
	}
}

func TestMarshalUnmarshalBinary(t *testing.T) {
	type TestBinary struct {
		Byte      byte
		private   byte
		IgnoreTag byte `binary:"-"`
		U8        uint8
		I8        int8
		U16       uint16
		I16       int16
		Bool      bool
	}

	b := TestBinary{
		Byte:      0xa1,
		private:   0xbb,
		IgnoreTag: 0xee,
		U8:        255,
		I8:        -128,
		U16:       65535,
		I16:       -32768,
		Bool:      true,
	}

	mb, err := Marshal(ProtocolTypeBINARY, b)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var check TestBinary
	err = Unmarshal(mb, &check)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if check.Byte != b.Byte ||
		check.private != 0x00 || // should ignore this tag
		check.IgnoreTag != 0x00 || // should ignore this tag
		check.U8 != b.U8 ||
		check.I8 != b.I8 ||
		check.U16 != b.U16 ||
		check.I16 != b.I16 ||
		check.Bool != b.Bool {
		t.Errorf("want: %#v. got: %#v\n", b, check)
	}

}
