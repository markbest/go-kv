package app

import (
	"net"

	"github.com/markbest/go-kv/utils"
	"github.com/markbest/go-kv/utils/tcp"
)

// send msg and receive reply
func senderMsgAndReceive(conn net.Conn, msg []byte) (string, error) {
	conn.Write(utils.Packet(msg))
	resp := make(chan []byte)

	go func() {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err == nil {
			resp <- buffer
		}

	}()
	return string(<-resp), nil
}

// handle send server msg
func HandleSendServerMsg(conn net.Conn, action, inputStr1, inputStr2 string) (string, error) {
	switch action {
	case "set":
		msg := "+:" + inputStr1 + ":" + inputStr2
		return senderMsgAndReceive(conn, []byte(msg))
	case "get":
		msg := "g:" + inputStr1
		return senderMsgAndReceive(conn, []byte(msg))
	case "delete":
		msg := "-:" + inputStr1 + ":" + inputStr2
		return senderMsgAndReceive(conn, []byte(msg))
	case "list":
		msg := "l:" + inputStr1
		return senderMsgAndReceive(conn, []byte(msg))
	case "persistent":
		msg := "p"
		return senderMsgAndReceive(conn, []byte(msg))
	case "exit":
		msg := "e"
		return senderMsgAndReceive(conn, []byte(msg))
	}
	return "", nil
}

// clear scan string
func ClearScan(inputStr1, inputStr2 *string) {
	*inputStr1 = ""
	*inputStr2 = ""
}

// handle scan input
func HandleScanInput(client *tcp.TCPClient, action, inputStr1, inputStr2 string) string {
	rs, err := client.ReadWrite(func(conn *net.TCPConn) (string, error) {
		return HandleSendServerMsg(conn, action, inputStr1, inputStr2)
	})
	if err != nil {
		return err.Error()
	} else {
		return rs
	}
}
