package main

import (
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/config"
	"github.com/6qhtsk/sonolus-test-server/controller"
	"github.com/6qhtsk/sonolus-test-server/service"
	"github.com/6qhtsk/sonolusgo"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Sonolus-Test-Server %s", config.Version)
	sonolusConfig := sonolusgo.DefaultConfig()
	sonolusConfig.ServerName = fmt.Sprintf("Ayachan测试服务器 - %s", config.Version)
	sonolusConfig.ServerBanner = sonolusgo.NewSRLServerBanner("daae4b4a3d9fe51bd76ab68457ce1e3c0443f39a", "/sonolus/repository/BackgroundImage/daae4b4a3d9fe51bd76ab68457ce1e3c0443f39a")
	sonolusConfig.Handlers.Levels = service.LevelHandlers
	router := gin.Default()
	sonolusConfig.LoadHandlers(router)
	sonolusConfig.RouterGroups.Levels.GET("/:name/bgm.mp3", controller.GetLevelItems("bgm.mp3"))
	sonolusConfig.RouterGroups.Levels.GET("/:name/data", controller.GetLevelItems("data"))
	router.POST("/upload", controller.ChartUploadHandler())
	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}
	return
}
