package app

import (
	"github.com/jmoiron/sqlx"
)

func Create(db *sqlx.DB) {
	_, err := db.Exec(`create table user (
		id integer primary key autoincrement,
		username string not null,
		password string not null,
		display_name string)`)
	if err != nil {
		panic(err)
	}
}
