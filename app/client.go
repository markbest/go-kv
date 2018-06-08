package app

import (
	"encoding/json"
	"net"

	"github.com/markbest/go-kv/utils"
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
		msg := &Msg{Action: "set", Key: inputStr1, Value: inputStr2}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return "", err
		}
		return senderMsgAndReceive(conn, jsonMsg)
	case "get":
		msg := &Msg{Action: "get", Key: inputStr1}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return "", err
		}
		return senderMsgAndReceive(conn, jsonMsg)
	case "delete":
		msg := &Msg{Action: "delete", Key: inputStr1}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return "", err
		}
		return senderMsgAndReceive(conn, jsonMsg)
	case "list":
		msg := &Msg{Action: "list"}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return "", err
		}
		return senderMsgAndReceive(conn, jsonMsg)
	case "persistent":
		msg := &Msg{Action: "persistent"}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return "", err
		}
		return senderMsgAndReceive(conn, jsonMsg)
	case "exit":
		msg := &Msg{Action: "exit"}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return "", err
		}
		return senderMsgAndReceive(conn, jsonMsg)
	}
	return "", nil
}
