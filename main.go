package main

import (
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/config"
	"github.com/6qhtsk/sonolus-test-server/controller"
	"github.com/6qhtsk/sonolus-test-server/service"
	"github.com/6qhtsk/sonolusgo"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Printf("Sonolus-Test-Server %s", config.Version)
	sonolusConfig := sonolusgo.DefaultConfig()
	sonolusConfig.ServerName = fmt.Sprintf("Ayachan测试服务器 - %s", config.Version)
	sonolusConfig.Handlers.Levels = service.LevelHandlers
	sonolusConfig.ServerBanner = sonolusgo.NewSRLServerBanner("daae4b4a3d9fe51bd76ab68457ce1e3c0443f39a", "https://repository.ayachan.fun/sonolus/BackgroundImage/daae4b4a3d9fe51bd76ab68457ce1e3c0443f39a")
	router := gin.Default()
	sonolusConfig.LoadHandlers(router)
	if !config.ServerCfg.UseTencentCos {
		sonolusConfig.RouterGroups.Levels.GET("/:name/bgm", controller.GetLevelItems("bgm"))
		sonolusConfig.RouterGroups.Levels.GET("/:name/data", controller.GetLevelItems("data"))
	}
	sonolusConfig.RouterGroups.Levels.POST("", controller.ChartUploadHandler())
	err := router.Run(":" + config.ServerCfg.Port)
	if err != nil {
		panic(err)
	}
	return
}
