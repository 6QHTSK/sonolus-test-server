package service

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
        CREATE TABLE post (id INTEGER PRIMARY KEY, title TEXT, difficulty INTEGER, hidden INTEGER, expired TIMESTAMP, bgmHash TEXT, dataHash TEXT);
        `
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func initDatabase() {
	dbFile := "./sonolus/database.db"
	err := initializeDatabase(dbFile)
	if err != nil {
		panic(err)
	}

	db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
}

func generatePostUid() int {
	var id int
	query := `
        WITH RECURSIVE cnt(x) AS (
            SELECT 100000
            UNION ALL
            SELECT x+1 FROM cnt
            LIMIT 900000
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

func insertPost(uid int, post model.UploadPost, bgmHash string, dataHash string) error {
	current := time.Now().UTC().Unix()
	expired := current + post.Lifetime
	_, err := db.Exec(`INSERT INTO post VALUES (?,?,?,?,?,?,?)`, uid, post.Title, post.Difficulty, post.Hidden, expired, bgmHash, dataHash)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func deleteDBOutdatedPost() (deleteUidList []int, err error) {
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
		err = db.Select(&postList, `SELECT * from post where hidden=FALSE LIMIT 20 OFFSET ?`, offset)
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
