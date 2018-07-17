package main

import (
	"flag"
	"os"

	"github.com/markbest/go-kv/app"
	"github.com/markbest/go-kv/conf"
	"github.com/markbest/go-kv/utils"
	"github.com/markbest/go-kv/utils/tcp"

	log "github.com/sirupsen/logrus"
)

var (
	kvConnect      *utils.KV
	configFilePath = flag.String("c", "env.yaml", "config file path")
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	// parse env.yaml
	flag.Parse()
	err := conf.ParseConfig(*configFilePath)
	if err != nil {
		panic(err)
	}

	// init kv data
	log.Info("init kv data from db")
	kvConnect = utils.NewKV(conf.Config.KV.DBPath)
	if count := kvConnect.Init(); count > 0 {
		log.Infof("restore %d data from db", count)
	}

	// persistent data
	if conf.Config.KV.Persistent {
		go app.PersistentTicker(kvConnect)
	}

	// start server
	server := tcp.NewTCPServer(conf.Config.App.Addr, conf.Config.App.Port, 10)
	lis, lisErr := server.Listen()
	if lisErr != nil {
		log.Errorf("start server occur error: %s", lisErr.Error())
	}
	log.Infof("success start server - %s:%s", conf.Config.App.Addr, conf.Config.App.Port)

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go app.HandleClientConnection(conn, kvConnect)
	}
}
