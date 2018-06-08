package tcp

import (
	"net"
)

type TCPClient struct {
	addr     string
	port     string
	conn     *net.TCPConn
	maxRetry int
}

// new tcp client
func NewTCPClient(addr, port string, maxRetry int) *TCPClient {
	return &TCPClient{
		addr:     addr,
		port:     port,
		conn:     nil,
		maxRetry: maxRetry,
	}
}

// tcp connect
func (t *TCPClient) connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", t.addr+":"+t.port)
	if err != nil {
		return err
	}

	var i = 0
	for {
		conn, connErr := net.DialTCP("tcp", nil, tcpAddr)
		if connErr == nil && conn != nil {
			t.conn = conn
			break
		}

		if i > t.maxRetry {
			return connErr
		}
		i += 1
	}
	return nil
}

// use tcp to write and read
func (t *TCPClient) ReadWrite(rw func(conn *net.TCPConn) (string, error)) (string, error) {
	if connErr := t.connect(); connErr != nil {
		return "", connErr
	}
	return rw(t.conn)
}

// close tcp connect
func (t *TCPClient) Close() error {
	if t.conn == nil {
		return nil
	}

	if closeErr := t.conn.Close(); closeErr != nil {
		return closeErr
	}
	t.conn = nil
	return nil
}
