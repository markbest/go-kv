package app

import (
	"net"
	"strings"
	"time"

	"github.com/markbest/go-kv/conf"
	"github.com/markbest/go-kv/utils"

	log "github.com/sirupsen/logrus"
)

// handle receive msg
func handleReceiveMsg(msg []byte, kv *utils.KV) (string, string, error) {
	receiveMsg := strings.Split(string(msg), ":")

	switch receiveMsg[0] {
	case "+":
		kv.Set(receiveMsg[1], receiveMsg[2])
		return "+", receiveMsg[2], nil
	case "g":
		value, err := kv.Get(receiveMsg[1])
		return "-", value, err
	case "-":
		kv.Del(receiveMsg[1])
		return "-", receiveMsg[1], nil
	case "l":
		keys := kv.List()
		if len(keys) > 0 {
			return "l", strings.Join(keys, "\n"), nil
		} else {
			return "l", "null", nil
		}
	case "p":
		kv.Persistent()
		return "p", "success persistent", nil
	case "e":
		return "e", "success disconnect", nil
	}
	return "", "", nil
}

// persistent data
func PersistentTicker(kv *utils.KV) {
	ticker := time.NewTicker(time.Duration(conf.Config.KV.PersistentTime) * time.Second)
	for {
		select {
		case <-ticker.C:
			kv.Persistent()
			log.Infof("kv persistent success")
		}
	}
}

// handle client request
func HandleClientConnection(conn net.Conn, kv *utils.KV) {
	tmpBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 16)
	go handleRequest(conn, readerChannel, kv)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		tmpBuffer = utils.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}

// handle receive request
func handleRequest(conn net.Conn, readerChannel chan []byte, kv *utils.KV) {
	for {
		select {
		case data := <-readerChannel:
			log.Infof("receive msg: %s", data)
			action, rs, err := handleReceiveMsg(data, kv)
			if action == "e" {
				conn.Write([]byte(rs))
				log.Infof("reply msg: %s", rs)
				conn.Close()
				continue
			}
			if err != nil {
				conn.Write([]byte(err.Error()))
				log.Infof("reply msg: %s", err.Error())
				continue
			}
			conn.Write([]byte(rs))
			log.Infof("reply msg: %s", rs)
		}
	}
}
