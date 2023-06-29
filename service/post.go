package service

import (
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/errors"
	"github.com/6qhtsk/sonolus-test-server/model"
	"github.com/h2non/filetype"
	"io"
	"log"
	"os"
	"sync"
)

var saveMutex sync.Mutex

func SavePost(post model.UploadPost) (uid int, err error, myError *errors.TestServerError) {
	saveMutex.Lock()
	defer saveMutex.Unlock()
	deleteOutdatedPost()
	uid = generatePostUid()
	err = createPostDir(uid)
	if err != nil {
		return 0, err, errors.FailCreateDir
	}

	// 处理保存音频文件任务
	src, err := post.Bgm.Open()
	if err != nil {
		return 0, err, errors.BGMProcessError
	}
	defer src.Close()
	// 检查音频文件是否正确
	head := make([]byte, 261)
	_, err = src.Read(head)
	if err != nil {
		return 0, err, errors.BGMProcessError
	}
	if !filetype.IsAudio(head) {
		trueFileType, _ := filetype.Match(head)
		return 0, fmt.Errorf("the file you upload is %s (MIME %s), not audio",
			trueFileType.Extension, trueFileType.MIME.Value), errors.BadBGMType
	}
	_, err = src.Seek(0, 0)
	if err != nil {
		return 0, err, errors.BGMProcessError
	}
	dst := getBgmPath(uid)
	// 创建音频文件
	out, err := os.Create(dst)
	if err != nil {
		return uid, err, errors.BGMProcessError
	}
	defer out.Close()
	// 保存文件
	_, err = io.Copy(out, src)
	if err != nil {
		return uid, err, errors.FailCreateFile
	}

	sonolusChart, err := post.Chart.ConvertToSonnolus()
	if err != nil {
		return 0, err, errors.ChartConvertFail
	}
	err = writeSonolusData(sonolusChart, getDataPath(uid))
	if err != nil {
		return uid, err, errors.FailCreateFile
	}
	err = insertPost(uid, post)
	if err != nil {
		return 0, err, errors.FailInsertDatabase
	}
	return uid, nil, nil
}

func deleteOutdatedPost() {
	outdatedPost, err := deleteDBOutdatedPost()
	if err != nil {
		log.Printf("删除谱面数据库条目时发生错误：%s", err)
		return
	}
	err = removeOutdatedPostDir(outdatedPost)
	if err != nil {
		log.Printf("删除谱面字段时发生错误：%s", err)
		return
	}
}
