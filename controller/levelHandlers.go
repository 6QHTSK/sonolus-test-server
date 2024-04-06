package controller

import (
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/config"
	"github.com/6qhtsk/sonolus-test-server/errors"
	"github.com/6qhtsk/sonolus-test-server/manager"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetLevelItems(itemName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("name")
		id, err := strconv.Atoi(idStr)
		if abortWhenErr(ctx, err, errors.BadUidErr) {
			return
		}
		ctx.File(fmt.Sprintf("./sonolus/levels/%d.%s", id, map[string]string{"bgm": "mp3", "data": "json.gz", "bdv2": "bdv2.json"}[itemName]))
	}
}

func RedirectBDV2Chart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("name")
		id, err := strconv.Atoi(idStr)
		if abortWhenErr(ctx, err, errors.BadUidErr) {
			return
		}
		remoteUrl := fmt.Sprintf("%s/%s", config.ServerCfg.Cos.AccessUrl, manager.GetCosBDV2DataPath(id))
		ctx.Redirect(http.StatusMovedPermanently, remoteUrl)
	}
}
