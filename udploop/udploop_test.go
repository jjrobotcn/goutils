package udploop

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestClientReaderServerWriter(t *testing.T) {
	cr, err := NewClientReader()
	if err != nil {
		assert.NoError(t, err)
	}

	sw, err := NewServerWriter()
	if err != nil {
		assert.NoError(t, err)
	}

	check := []byte{1, 3, 5, 7, 9}

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()

		buf := make([]byte, 1024)

		n, err := cr.Read(buf)
		if err != nil {
			assert.NoError(t, err)
		}

		got := buf[:n]

		assert.Equal(t, got, check)
	}()

	if _, err := sw.Write(check); err != nil {
		t.Error(err)
	}

	wg.Wait()
}

func TestServerReaderClientWriter(t *testing.T) {
	sr, err := NewServerReader()
	if err != nil {
		assert.NoError(t, err)
	}

	cw, err := NewClientWriter()
	if err != nil {
		assert.NoError(t, err)
	}

	check := []byte{1, 3, 5, 7, 9}

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()

		buf := make([]byte, 1024)

		n, err := sr.Read(buf)
		if err != nil {
			assert.NoError(t, err)
		}

		got := buf[:n]

		assert.Equal(t, got, check)
	}()

	if _, err := cw.Write(check); err != nil {
		assert.NoError(t, err)
	}

	wg.Wait()
}
