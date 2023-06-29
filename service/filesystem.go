package service

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/model"
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
	return fmt.Sprintf("%s/%d/data", localRepo, uid)
}

func getBgmPath(uid int) string {
	return fmt.Sprintf("%s/%d/bgm", localRepo, uid)
}

func createPostDir(uid int) error {
	return os.MkdirAll(fmt.Sprintf("%s/%d", localRepo, uid), os.FileMode(0755))
}

func writeSonolusData(chart model.SonolusLevelData, dest string) error {
	data, err := json.Marshal(chart)
	if err != nil {
		return err
	}

	dataFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	gw := gzip.NewWriter(dataFile)
	defer gw.Close()
	_, err = gw.Write(data)
	if err != nil {
		return err
	}
	return err
}

func removeOutdatedPostDir(outdatedPostUid []int) error {
	for _, uid := range outdatedPostUid {
		err := os.RemoveAll(fmt.Sprintf("%s/%d", localRepo, uid))
		if err != nil {
			return err
		}
	}
	return nil
}
