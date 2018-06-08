package tcp

import "net"

type TCPServer struct {
	addr     string
	port     string
	maxRetry int
}

// new tcp server
func NewTCPServer(addr, port string, maxRetry int) *TCPServer {
	return &TCPServer{
		addr:     addr,
		port:     port,
		maxRetry: maxRetry,
	}
}

// server listen
func (t *TCPServer) Listen() (lis net.Listener, lisErr error) {
	for i := 0; i < t.maxRetry; i++ {
		lis, lisErr = net.Listen("tcp", t.addr+":"+t.port)
		if lisErr == nil {
			break
		}
	}
	return lis, lisErr
}
