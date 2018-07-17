package conf

import (
	"github.com/toolkits/file"
	"gopkg.in/yaml.v2"
)

var Config *GlobalConfig

type AppConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
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
		panic(err.Error())
	}

	var c GlobalConfig
	err = yaml.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return err
	}
	Config = &c
	return nil
}
