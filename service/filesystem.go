package service

import (
	"bytes"
	"fmt"
	"github.com/h2non/filetype"
	"io"
	"log"
	"os"
)

var localRepo = "./sonolus/levels"

func initLocalRepo() {
	err := os.MkdirAll(localRepo, os.FileMode(0755))
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func getDataPath(uid int) string {
	return fmt.Sprintf("%s/%d.json.gz", localRepo, uid)
}

func getBgmPath(uid int) string {
	return fmt.Sprintf("%s/%d.mp3", localRepo, uid)
}

func getBDV2DataPath(uid int) string {
	return fmt.Sprintf("%s/%d.bdv2.json", localRepo, uid)
}

func removeOutdatedPost(outdatedPostUid []int) error {
	for _, uid := range outdatedPostUid {
		for _, item := range []string{getDataPath(uid), getBgmPath(uid), getBDV2DataPath(uid)} {
			err := os.Remove(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkIfAudio(bgmData []byte) error {
	head := bgmData[0:261]
	if !filetype.IsAudio(head) {
		trueFileType, _ := filetype.Match(head)
		return fmt.Errorf("the file you upload is %s (MIME %s), not audio",
			trueFileType.Extension, trueFileType.MIME.Value)
	}
	return nil
}

func saveIOReaderToFile(file io.Reader, dest string) error {
	// 创建文件
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	// 保存文件
	_, err = io.Copy(out, file)
	return err
}

func saveBytesToFile(data []byte, dest string) (err error) {
	reader := bytes.NewReader(data)
	return saveIOReaderToFile(reader, dest)
}
