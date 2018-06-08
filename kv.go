package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/markbest/go-kv/app"
	"github.com/markbest/go-kv/conf"
	"github.com/markbest/go-kv/utils"
	"github.com/markbest/go-kv/utils/tcp"

	log "github.com/sirupsen/logrus"
)

var (
	kvConnect      *utils.KV
	configFilePath = flag.String("c", "go-kv/env.yaml", "config file path")
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
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

	// pprof server
	if conf.Config.App.Debug {
		log.Infof("start pprof server: %s", conf.Config.App.Pprof)
		pprofServer := &http.Server{Addr: conf.Config.App.Pprof}
		go pprofServer.ListenAndServe()
	}

	// persistent data
	if conf.Config.KV.Persistent {
		go app.PersistentTicker(kvConnect)
	}

	// start server
	server := tcp.NewTCPServer(conf.Config.App.ListenAddr, conf.Config.App.ListenPort, 10)
	lis, lisErr := server.Listen()
	if lisErr != nil {
		log.Errorf("start server occur error: %s", lisErr.Error())
	}
	log.Infof("success start server - %s:%s", conf.Config.App.ListenAddr, conf.Config.App.ListenPort)

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go app.HandleClientConnection(conn, kvConnect)
	}
}
