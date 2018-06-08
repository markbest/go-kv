package conf

import (
	"github.com/toolkits/file"
	"gopkg.in/yaml.v2"
	"log"
)

var Config *GlobalConfig

type AppConfig struct {
	ListenAddr string `yaml:"listenaddr"`
	ListenPort string `yaml:"listenport"`
	Debug      bool   `yaml:"debug"`
	Pprof      string `yaml:"pprof"`
}

type KVConfig struct {
	DBPath         string `yaml:"dbpath"`
	Persistent     bool   `yaml:"persistent"`
	PersistentTime int64  `yaml:"persistenttime"`
}

type GlobalConfig struct {
	App *AppConfig `yaml:"app"`
	KV  *KVConfig  `yaml:"kv"`
}

func ParseConfig(cfg string) error {
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return err
		log.Fatalf("read config file error: %s", err.Error())
	}

	var c GlobalConfig
	err = yaml.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return err
	}
	Config = &c
	return nil
}
