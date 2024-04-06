package dao

import "github.com/jmoiron/sqlx"

func init() {
	dbFile := "./sonolus/database.db"
	err := initializeDatabase(dbFile)
	if err != nil {
		panic(err)
	}

	db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	addUploadTimeColumn()
}
