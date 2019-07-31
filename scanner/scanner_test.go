package scanner

import (
	"bufio"
	"bytes"
	"testing"
)

func TestSerialSplit(t *testing.T) {
	ps := [][]byte{
		{0xf0},
		{0, 0xf0},
		{0x0f, 1, 0xf0},
		{0x0e, 0xf0},
		{0x0f, 2, 0xf0},
		{0x0e, 0xf0},
		{0x0f, 3, 0xf0},
		{0x0e, 0xf0},
		{0x0f, 4, 0xf0},
		{0x0f},
	}

	multi := new(bytes.Buffer)
	for _, p := range ps {
		multi.Write(p)
	}

	s := bufio.NewScanner(multi)
	s.Split(SerialSplit)

	count := 0
	for s.Scan() {
		if !bytes.Equal(ps[count], s.Bytes()) {
			t.Errorf("want: %#v, got: %#v", ps[count], s.Bytes())
		}
		count++
	}
	if count != 10 {
		t.Errorf("want split count 10, got: %d", count)
	}
}
