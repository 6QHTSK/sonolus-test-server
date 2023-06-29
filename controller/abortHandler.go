package controller

import (
	"github.com/6qhtsk/sonolus-test-server/errors"
	"github.com/gin-gonic/gin"
	"log"
)

func abortWhenErr(ctx *gin.Context, err error, myError *errors.TestServerError) bool {
	if err != nil {
		log.Printf("%s : %s", myError, err)
		ctx.JSON(myError.HttpCode, gin.H{
			"code":        myError.ErrCode,
			"description": myError.ErrMsg,
			"detail":      err.Error(),
		})
		ctx.Abort()
		return true
	}
	return false
}
