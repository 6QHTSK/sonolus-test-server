package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/config"
	"github.com/6qhtsk/sonolus-test-server/errors"
	"github.com/6qhtsk/sonolus-test-server/model"
	"io"
	"log"
	"sync"
)

var saveMutex sync.Mutex

func SavePost(post model.UploadPost) (uid int, err error, myError *errors.TestServerError) {
	saveMutex.Lock()
	defer saveMutex.Unlock()
	deleteOutdatedPost()
	uid = generatePostUid()

	// 处理保存音频文件任务
	// 音频文件大于20M时，返回错误
	if post.Bgm.Size >= 20*1024*1024 {
		return 0, fmt.Errorf("bgm too big >20MB ( %.1f MB)", float64(post.Bgm.Size)/1024.0/1024.0), errors.FileTooBig
	}
	if err != nil {
		return 0, err, errors.ChartConvertFail
	}
	// 打开音频文件，检查音频文件是否正确
	bgmFile, err := post.Bgm.Open()
	if err != nil {
		return 0, err, errors.BGMProcessError
	}
	defer bgmFile.Close()
	bgmBuffer := bytes.NewBuffer(nil)
	_, err = io.Copy(bgmBuffer, bgmFile)
	if err != nil {
		return 0, err, errors.BGMProcessError
	}
	err = checkIfAudio(bgmBuffer.Bytes())
	if err != nil {
		return 0, err, errors.BadBGMType
	}
	bgmHash := bytesSha1(bgmBuffer.Bytes())
	// 转码Sonolus谱面
	sonolusChart, err := post.Chart.ConvertToSonnolus()
	if err != nil {
		return 0, err, errors.ChartConvertFail
	}
	sonolusChartRawData, err := json.Marshal(sonolusChart)
	if err != nil {
		return 0, err, errors.ChartConvertFail
	}
	sonolusChartData, err := gzippedBytes(sonolusChartRawData)
	if err != nil {
		return 0, err, errors.ChartConvertFail
	}
	datahash := bytesSha1(sonolusChartData)
	// 保存文件
	if config.ServerCfg.UseTencentCos { // 保存到腾讯云
		err = uploadBytesToTencentCos(bgmBuffer.Bytes(), getCosBgmPath(uid))
		if err != nil {
			return 0, err, errors.FailUploadToTencentCos
		}
		err = uploadBytesToTencentCos(sonolusChartData, getCosDataPath(uid))
		if err != nil {
			return 0, err, errors.FailUploadToTencentCos
		}
	} else { // 保存到本地系统
		err = saveBytesToFile(bgmBuffer.Bytes(), getBgmPath(uid))
		if err != nil {
			return 0, err, errors.FailCreateFile
		}
		err = saveBytesToFile(sonolusChartData, getDataPath(uid))
		if err != nil {
			return uid, err, errors.FailCreateFile
		}
	}
	// 插入到数据库
	if insertPost(uid, post, bgmHash, datahash) != nil {
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
	if config.ServerCfg.UseTencentCos {
		for _, uid := range outdatedPost {
			err = deleteInTencentCos([]string{getCosDataPath(uid), getCosBgmPath(uid)})
			if err != nil {
				log.Printf("删除谱面字段时发生错误：%s", err)
				return
			}
		}
	} else {
		err = removeOutdatedPost(outdatedPost)
		if err != nil {
			log.Printf("删除谱面字段时发生错误：%s", err)
			return
		}
	}
}
