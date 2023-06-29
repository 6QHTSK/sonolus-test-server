package model

import (
	"encoding/json"
	"mime/multipart"
)

type UploadPost struct {
	Title      string                `form:"title" binding:"required"`
	Bgm        *multipart.FileHeader `form:"bgm" binding:"required"`
	Chart      BestdoriChart
	ChartStr   string `form:"chart" binding:"required"`
	Difficulty int    `form:"difficulty"`
	Hidden     bool   `form:"hidden"`
	Lifetime   int64  `form:"lifetime"`
}

func (post *UploadPost) ParseChart() error {
	return json.Unmarshal([]byte(post.ChartStr), &post.Chart)
}
