package database

import (
	"database/sql"
	"os"

	//sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	//DBNAME is a constant database filename.
	DBNAME string = "sqlite-database.db"
)

var (
	//DBConn is a SQLite3 database connection.
	DBConn *sql.DB
)

//InitDatabase initialises a SQLite3  database.
func InitDatabase() error {
	//If database file doesn't exist, create it.
	if !fileExists(DBNAME) {
		file, err := os.Create(DBNAME)
		if err != nil {
			return err
		}

		file.Close()
	}

	db, err := sql.Open("sqlite3", "./"+DBNAME)
	if err != nil {
		return err
	}

	err = prepare(db)
	if err != nil {
		return err
	}

	DBConn = db
	return nil
}

func prepare(db *sql.DB) error {
	/*
		Prepares SQLite database.
		First query enables ON DELETE constrains.
		Rest create models and unique indexes.
	*/

	queries := []string{
		`PRAGMA foreign_keys=ON
	`, `
	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY,
		body TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	`, `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_questions_body 
    ON questions (body);
	`, `
	CREATE TABLE IF NOT EXISTS options (
		id INTEGER PRIMARY KEY,
		question_id INTEGER NOT NULL,
		body TEXT NOT NULL,
		correct BOOLEAN NOT NULL DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
	);
	`, `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_options_body 
    ON options (body);
	`}

	for _, query := range queries {
		statement, err := db.Prepare(query)
		if err != nil {
			return err
		}

		_, err = statement.Exec()
		if err != nil {
			return err
		}

		statement.Close()
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
