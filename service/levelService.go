package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/6qhtsk/sonolus-test-server/errors"
	"github.com/6qhtsk/sonolus-test-server/model"
	"github.com/6qhtsk/sonolusgo"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func getNoticeLevel() sonolusgo.Level {
	engine, _, err := sonolusgo.GetItem[sonolusgo.Engine]("bandori")
	if err != nil {
		panic(err)
	}
	return sonolusgo.Level{
		Name:          "置顶公告#",
		Version:       1,
		Rating:        35,
		Title:         "谱面上传地址：https://www.ayachan.fun/#/sonolus-upload\n隐藏谱面请通过搜索界面输入隐藏谱面UID游玩",
		Artists:       "Beta测试中，不稳定\n遇到各种问题可以通过QQ联系我或提Github Issue",
		Author:        "测试服务器作者：彩绫与6QHTSK",
		Engine:        engine,
		UseSkin:       sonolusgo.UseItem[sonolusgo.Skin]{UseDefault: true},
		UseBackground: sonolusgo.UseItem[sonolusgo.Background]{UseDefault: true},
		UseEffect:     sonolusgo.UseItem[sonolusgo.Effect]{UseDefault: true},
		UseParticle:   sonolusgo.UseItem[sonolusgo.Particle]{UseDefault: true},
		Cover: sonolusgo.SRLLevelCover{
			Type: "LevelCover",
			Hash: fileSha1("./sonolus/repository/LevelThumbnail/notice"),
			Url:  "/sonolus/repository/LevelThumbnail/notice",
		},
		Bgm:     sonolusgo.SRLLevelBgm{},
		Preview: nil,
		Data:    sonolusgo.SRLLevelData{},
	}
}

var LevelHandlers = sonolusgo.SonolusHandlers[sonolusgo.Level]{
	List:      LevelListHandler,
	Search:    LevelSearchHandler,
	Item:      LevelItemHandler,
	Recommend: sonolusgo.GetEmptyRecommend[sonolusgo.Level],
}

func fileSha1(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func LevelListHandler(page int, queryMap map[string]string) (pageCount int, items []sonolusgo.Level) {
	name, ok := queryMap["uid"]
	var uid int
	if name == "" || !ok { // 未指定uid
		uid = -1
	} else {
		var err error
		uid, err = strconv.Atoi(name)
		if err != nil {
			return 0, []sonolusgo.Level{}
		}
	}
	deleteOutdatedPost()
	postCnt, err := GetPostCnt(uid)
	if err != nil {
		return 0, []sonolusgo.Level{}
	}
	pageCount = postCnt/20 + 1
	posts, err := GetPost(uid, page*20)
	if err != nil {
		return 0, []sonolusgo.Level{}
	}
	if page == 0 && uid == -1 {
		items = append(items, getNoticeLevel())
	}
	for _, post := range posts {
		items = append(items, convertDatabaseToSonolus(post))
	}
	return pageCount, items
}

func LevelSearchHandler() (search sonolusgo.Search) {
	search.Options = append(search.Options, sonolusgo.NewSearchTextOption("uid", "UID - 可搜索隐藏谱面", "隐藏谱面ID"))
	return search
}

func LevelItemHandler(name string) (item sonolusgo.Level, description string, err error) {
	if name == "置顶公告#" {
		return getNoticeLevel(), "公告不可游玩(´・ω・`)\n可在Sonolus1群、BanGDream谱师群或SonolusQQ频道找到我\nGithub项目：还没上传，记得上传后在这里改\nB站：@彩绫与6QHTSK 关私信提醒了你的消息可能很久才会被我看到", nil
	}
	uid, err := strconv.Atoi(name)
	if err != nil {
		return item, description, err
	}
	deleteOutdatedPost()
	post, err := GetPost(uid, 0)
	if len(post) != 1 {
		return item, description, errors.BadUidErr
	}
	item = convertDatabaseToSonolus(post[0])
	return item, "", nil
}

func convertDatabaseToSonolus(dbItem model.DatabasePost) sonolusgo.Level {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	engine, _, err := sonolusgo.GetItem[sonolusgo.Engine]("bandori")
	if err != nil {
		panic(err)
	}
	r := rand.New(rand.NewSource(int64(dbItem.Id)))
	randCoverID := r.Intn(35) + 1
	return sonolusgo.Level{
		Name:          strconv.Itoa(dbItem.Id),
		Version:       1,
		Rating:        dbItem.Difficulty,
		Title:         dbItem.Title,
		Artists:       fmt.Sprintf("Expired at %s", dbItem.Expired.In(loc).Format(time.DateTime)),
		Author:        map[bool]string{false: "一般测试谱面", true: "隐藏测试谱面"}[dbItem.Hidden],
		Engine:        engine,
		UseSkin:       sonolusgo.UseItem[sonolusgo.Skin]{UseDefault: true},
		UseBackground: sonolusgo.UseItem[sonolusgo.Background]{UseDefault: true},
		UseEffect:     sonolusgo.UseItem[sonolusgo.Effect]{UseDefault: true},
		UseParticle:   sonolusgo.UseItem[sonolusgo.Particle]{UseDefault: true},
		Cover: sonolusgo.SRLLevelCover{
			Type: "LevelCover",
			Hash: fileSha1(fmt.Sprintf("./sonolus/repository/LevelThumbnail/%d", randCoverID)),
			Url:  fmt.Sprintf("/sonolus/repository/LevelThumbnail/%d", randCoverID),
		},
		Bgm: sonolusgo.SRLLevelBgm{
			Type: "LevelBgm",
			Hash: fileSha1(fmt.Sprintf("./sonolus/levels/%d/bgm", dbItem.Id)),
			Url:  fmt.Sprintf("/sonolus/levels/%d/bgm", dbItem.Id),
		},
		Preview: nil,
		Data: sonolusgo.SRLLevelData{
			Type: "LevelData",
			Hash: fileSha1(fmt.Sprintf("./sonolus/levels/%d/data", dbItem.Id)),
			Url:  fmt.Sprintf("/sonolus/levels/%d/data", dbItem.Id),
		},
	}
}
