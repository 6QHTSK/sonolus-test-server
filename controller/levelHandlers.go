package controller

import (
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetLevelItems(itemName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("name")
		id, err := strconv.Atoi(idStr)
		if abortWhenErr(ctx, err, errors.BadUidErr) {
			return
		}
		ctx.File(fmt.Sprintf("./sonolus/levels/%d/%s", id, itemName))
	}
}
