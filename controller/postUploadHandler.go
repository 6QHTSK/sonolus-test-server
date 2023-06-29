package controller

import (
	"github.com/6qhtsk/sonolus-test-server/errors"
	"github.com/6qhtsk/sonolus-test-server/model"
	"github.com/6qhtsk/sonolus-test-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChartUploadHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uploadChart model.UploadPost
		err := ctx.ShouldBind(&uploadChart)
		if abortWhenErr(ctx, err, errors.UploadFormBindErr) {
			return
		}
		if uploadChart.Difficulty <= 0 {
			uploadChart.Difficulty = 25
		}
		if uploadChart.Lifetime <= 0 {
			uploadChart.Lifetime = 21600
		}
		err = uploadChart.ParseChart()
		if abortWhenErr(ctx, err, errors.UploadChartErr) {
			return
		}
		uid, err, myError := service.SavePost(uploadChart)
		if abortWhenErr(ctx, err, myError) {
			return
		}
		ctx.JSON(http.StatusOK, struct {
			Uid int `json:"uid"`
		}{Uid: uid})
	}
}
