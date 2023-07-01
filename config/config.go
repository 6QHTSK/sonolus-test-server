package config

import (
	"github.com/go-ini/ini"
)

type ServerConfig struct {
	Port          string           `ini:"port"`
	UseTencentCos bool             `ini:"use-tencent-cos"`
	Cos           TencentCosConfig `ini:"tencent-cos"`
}

var ServerCfg ServerConfig

func newServerConfig() ServerConfig {
	return ServerConfig{
		Port:          "8000",
		UseTencentCos: false,
		Cos:           TencentCosConfig{},
	}
}

type TencentCosConfig struct {
	CosUrl    string `ini:"cos-url"`
	SecretID  string `ini:"secret-id"`
	SecretKey string `ini:"secret-key"`
}

func init() {
	ServerCfg = newServerConfig()
	err := ini.MapTo(&ServerCfg, "config.ini")
	if err != nil {
		cfg := ini.Empty()
		ServerCfg = newServerConfig()
		err = ini.ReflectFrom(cfg, &ServerCfg)
		if err != nil {
			panic(err)
		}
		err = cfg.SaveTo("config.ini")
		if err != nil {
			panic(err)
		}
	}
}
