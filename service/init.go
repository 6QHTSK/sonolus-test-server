package service

import "github.com/6qhtsk/sonolus-test-server/config"

func init() {
	if config.ServerCfg.UseTencentCos {
		initTencentCos()
	} else {
		initLocalRepo()
	}
	initDatabase()
}
