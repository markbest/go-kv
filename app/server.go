package app

import (
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/markbest/go-kv/conf"
	"github.com/markbest/go-kv/utils"

	log "github.com/sirupsen/logrus"
)

// handle receive msg
func handleReceiveMsg(msg []byte, kv *utils.KV) (string, string, error) {
	receiveMsg := &Msg{}
	err := json.Unmarshal(msg, receiveMsg)
	if err != nil {
		return receiveMsg.Action, "", err
	}

	switch receiveMsg.Action {
	case "set":
		kv.Set(receiveMsg.Key, receiveMsg.Value)
		return receiveMsg.Action, receiveMsg.Value, nil
	case "get":
		value, err := kv.Get(receiveMsg.Key)
		return receiveMsg.Action, value, err
	case "delete":
		kv.Del(receiveMsg.Key)
		return receiveMsg.Action, receiveMsg.Key, nil
	case "list":
		keys := kv.List()
		if len(keys) > 0 {
			return receiveMsg.Action, strings.Join(keys, "\n"), nil
		} else {
			return receiveMsg.Action, "null", nil
		}
	case "persistent":
		kv.Persistent()
		return receiveMsg.Action, "success persistent", nil
	case "exit":
		return receiveMsg.Action, "success disconnect", nil
	}
	return receiveMsg.Action, "", nil
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
			if action == "exit" {
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
