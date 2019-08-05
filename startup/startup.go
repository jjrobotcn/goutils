package startup

import (
	"bufio"
	"context"
	"github.com/jjrobotcn/goutils/encoding/serial"
	_ "github.com/jjrobotcn/goutils/logrusenv"
	"github.com/jjrobotcn/goutils/scanner"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

// Timeout is start up protocol timeout params
var Timeout = time.Second

type Protocol struct {
	Func string `json:"fun"`
	Key  string `json:"key"`
}

func StartUp(r io.Reader, w io.Writer) error {
	logger := log.WithField("handler", "startup.StartUp")

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	// startup request protocol
	req := Protocol{
		Func: "update",
		Key:  "noUpdate",
	}
	upReq, err := serial.Marshal(serial.ProtocolTypeJSON, req)
	if err != nil {
		logger.WithError(err).Error("marshal json startup protocol error")
		return err
	}

	// error channel
	var errChan = make(chan error)
	defer close(errChan)

	// listen response
	go func() {
		err := listenStartUp(ctx, r)
		select {
		case <-ctx.Done():
			return
		default:
			errChan <- err
		}
	}()

	// write protocol data to serial
	_, err = w.Write(upReq)
	if err != nil {
		logger.WithError(err).Error("write json startup protocol error")
		return err
	}

	// block until context done or got response
	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errChan:
		logger.WithError(e).Debug("startup response")
		return e
	}
}

// listen start up success protocol with blocked
func listenStartUp(ctx context.Context, r io.Reader) error {
	s := bufio.NewScanner(r)
	s.Split(scanner.SerialSplit)

	var resp Protocol

	for s.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := serial.Unmarshal(s.Bytes(), &resp); err != nil {
			return err
		}

		// start up success
		if resp.Func == "update" && resp.Key == "startUpSuccess" {
			return nil
		}
	}

	return s.Err()
}
