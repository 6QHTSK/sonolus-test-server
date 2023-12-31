package model

import "time"

type DatabasePost struct {
	Id         int       `db:"id"`
	Title      string    `db:"title"`
	Difficulty int       `db:"difficulty"`
	Expired    time.Time `db:"expired"`
	Upload     time.Time `db:"upload"`
	Hidden     bool      `db:"hidden"`
	BgmHash    string    `db:"bgmHash"`
	DataHash   string    `db:"dataHash"`
}
