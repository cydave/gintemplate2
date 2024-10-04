package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() (*sql.DB, error) {
	d, err := sql.Open("sqlite3", "/tmp/data.db")
	if err != nil {
		return nil, err
	}
	err = d.Ping()
	if err != nil {
		return nil, err
	}

	db = d
	return db, nil
}

func Migrate(DB *sql.DB) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// poor man's migrate logic; SQL statements here.
	/*	tx.Exec(`
		CREATE TABLE user(
			id VARCHAR(36) PRIMARY KEY,
			username VARCHAR(120) UNIQUE NOT NULL,
			password VARCHAR(256) NOT NULL,
			firstname VARCHAR(120),
			lastname VARCHAR(120),
			email VARCHAR(120)
		);
		`)
	*/

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func Get() *sql.DB {
	return db
}
