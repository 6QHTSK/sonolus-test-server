package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
)

var tencentCosClient *cos.Client

func initTencentCos() {
	u, err := url.Parse(config.ServerCfg.Cos.CosUrl)
	if err != nil {
		panic(err)
	}
	b := &cos.BaseURL{BucketURL: u}
	tencentCosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.ServerCfg.Cos.SecretID,
			SecretKey: config.ServerCfg.Cos.SecretKey,
		},
	})
}

const cosPathPrefix = "sonolus/test"

func getCosBgmPath(uid int) string {
	return fmt.Sprintf("%s/%d.mp3", cosPathPrefix, uid)
}

func getCosDataPath(uid int) string {
	return fmt.Sprintf("%s/%d.json.gz", cosPathPrefix, uid)
}

func uploadBytesToTencentCos(data []byte, filepath string) (err error) {
	reader := bytes.NewReader(data)
	return uploadToTencentCos(reader, filepath)
}

func uploadToTencentCos(file io.Reader, filepath string) (err error) {
	_, err = tencentCosClient.Object.Put(context.Background(), filepath, file, nil)
	return err
}

func deleteInTencentCos(filepaths []string) (err error) {
	var obs []cos.Object
	for _, v := range filepaths {
		obs = append(obs, cos.Object{Key: v})
	}
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
	}
	_, _, err = tencentCosClient.Object.DeleteMulti(context.Background(), opt)
	return err
}
