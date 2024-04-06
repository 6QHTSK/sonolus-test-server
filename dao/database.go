package dao

import (
	"github.com/6qhtsk/sonolus-test-server/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

var db *sqlx.DB

func initializeDatabase(dbFile string) error {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		db, err := sqlx.Open("sqlite3", dbFile)
		if err != nil {
			return err
		}
		defer db.Close()

		sqlStmt := `
        CREATE TABLE post (id INTEGER PRIMARY KEY, title TEXT, difficulty INTEGER, hidden INTEGER, expired TIMESTAMP, bgmHash TEXT, dataHash TEXT, upload TIMESTAMP);
        `
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// 临时，用于增加一列，用于记录上传时间
func addUploadTimeColumn() {
	sqlStmtAddUploadColumn := "ALTER TABLE post ADD COLUMN upload TIMESTAMP DEFAULT 0"
	sqlStmtChkUploadColumn := "SELECT count(*) from sqlite_master where name='post' and sql like '%upload%'"
	var result int
	err := db.Get(&result, sqlStmtChkUploadColumn)
	if err != nil {
		log.Fatal(err)
	}
	if result == 0 {
		_, err := db.Exec(sqlStmtAddUploadColumn)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GeneratePostUid() int {
	var id int
	query := `
        WITH RECURSIVE cnt(x) AS (
            SELECT 100
            UNION ALL
            SELECT x+1 FROM cnt
            LIMIT 900
        )
        SELECT x FROM cnt
        WHERE x NOT IN (SELECT id FROM post)
        ORDER BY RANDOM()
        LIMIT 1;
    `
	err := db.Get(&id, query)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func InsertPost(uid int, post model.UploadPost, bgmHash string, dataHash string) error {
	current := time.Now().UTC().Unix()
	expired := current + post.Lifetime
	_, err := db.Exec(`INSERT INTO post(id,title,difficulty,hidden,expired,bgmHash,dataHash,upload) VALUES (?,?,?,?,?,?,?,?)`, uid, post.Title, post.Difficulty, post.Hidden, expired, bgmHash, dataHash, current)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func DeleteDBOutdatedPost() (deleteUidList []int, err error) {
	current := time.Now().UTC().Unix()
	err = db.Select(&deleteUidList, `SELECT id from post where expired < ?`, current)
	if err != nil {
		return deleteUidList, err
	}
	_, err = db.Exec(`DELETE FROM POST where expired < ?`, current)
	if err != nil {
		return deleteUidList, err
	}
	return deleteUidList, nil
}

func GetPost(uid int, offset int) (postList []model.DatabasePost, err error) {
	if uid == -1 {
		err = db.Select(&postList, `SELECT * from post where hidden=FALSE ORDER BY upload DESC LIMIT 20 OFFSET ?`, offset)
	} else {
		err = db.Select(&postList, `SELECT * from post where id=? LIMIT 20 OFFSET ?`, uid, offset)
	}
	return postList, err
}

func GetPostCnt(uid int) (postCnt int, err error) {
	if uid == -1 {
		err = db.Get(&postCnt, `SELECT count(id) from post where hidden=FALSE`)
	} else {
		err = db.Get(&postCnt, `SELECT count(id)  from post where id=?`, uid)
	}
	return postCnt, err
}
