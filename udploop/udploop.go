package udploop

import (
	"io"
	"net"
)

var IP = net.IP{239, 0, 0, 0}

const (
	PortServerRecvFrom = 5000
	PortServerSendTo   = 6000
	PortClientRecvFrom = 6000
	PortClientSendTo   = 5000
)

// NewServerReader返回一个服务端Reader，用户端写入数据将由此Reader读取
func NewServerReader() (io.ReadCloser, error) {
	return newReader(PortServerRecvFrom)
}

// NewServerWriter返回一个服务端Writer，写入数据将由用户端Reader读取
func NewServerWriter() (io.WriteCloser, error) {
	return newWriter(PortServerSendTo)
}

// NewClientReader返回一个用户端Reader，服务端写入数据将由此Reader读取
func NewClientReader() (io.ReadCloser, error) {
	return newReader(PortClientRecvFrom)
}

// NewClientWriter返回一个用户端Writer，写入数据将由服务端Reader读取
func NewClientWriter() (io.WriteCloser, error) {
	return newWriter(PortClientSendTo)
}

func newWriter(port int) (io.WriteCloser, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}

	w := &writer{
		conn: conn,
		dst: &net.UDPAddr{
			IP:   IP,
			Port: port,
		},
	}

	return w, nil
}

type writer struct {
	conn *net.UDPConn
	dst  *net.UDPAddr
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.conn.WriteToUDP(p, w.dst)
}

func (w *writer) Close() error {
	return w.conn.Close()
}

func newReader(port int) (io.ReadCloser, error) {
	conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{IP: IP, Port: port})
	if err != nil {
		return nil, err
	}

	r := &reader{
		conn: conn,
	}

	return r, nil
}

type reader struct {
	conn *net.UDPConn
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, _, err = r.conn.ReadFromUDP(p)
	return
}

func (r *reader) Close() error {
	return r.conn.Close()
}
