package main

import (
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/config"
	"github.com/6qhtsk/sonolus-test-server/controller"
	"github.com/6qhtsk/sonolus-test-server/service"
	"github.com/6qhtsk/sonolusgo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Printf("Sonolus-Test-Server %s", config.Version)
	sonolusConfig := sonolusgo.DefaultConfig()
	sonolusConfig.ServerName = fmt.Sprintf("Ayachan测试服务器 - %s", config.Version)
	Banner := sonolusgo.SRL{
		Hash: "daae4b4a3d9fe51bd76ab68457ce1e3c0443f39a",
		Url:  "https://repository.ayachan.fun/sonolus/BackgroundImage/daae4b4a3d9fe51bd76ab68457ce1e3c0443f39a",
	}
	service.LevelHandlers.Banner = Banner
	sonolusConfig.Handlers.Levels = service.LevelHandlers
	sonolusConfig.ServerBanner = Banner
	if !config.ServerCfg.UseTencentCos {
		sonolusConfig.ServerInfo = sonolusgo.ServerInfoFilePath{
			Levels:      "",
			Skins:       "./sonolus/skins.local.json",
			Backgrounds: "./sonolus/backgrounds.local.json",
			Effects:     "./sonolus/effects.local.json",
			Particles:   "./sonolus/particles.local.json",
			Engines:     "./sonolus/engines.local.json",
		}
	}
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(controller.RequestSizeLimiter(30 * 1024 * 1024)) // limits all the uploads with in 30MB
	sonolusConfig.LoadHandlers(router)
	if !config.ServerCfg.UseTencentCos {
		sonolusConfig.RouterGroups.Levels.GET("/:name/bgm", controller.GetLevelItems("bgm"))
		sonolusConfig.RouterGroups.Levels.GET("/:name/data", controller.GetLevelItems("data"))
		sonolusConfig.RouterGroups.Levels.GET("/:name/bdv2.json", controller.GetLevelItems("bdv2"))
	} else {
		sonolusConfig.RouterGroups.Levels.GET("/:name/bdv2.json", controller.RedirectBDV2Chart())
	}
	sonolusConfig.RouterGroups.Levels.POST("", controller.ChartUploadHandler())
	err := router.Run(":" + config.ServerCfg.Port)
	if err != nil {
		panic(err)
	}
	return
}
