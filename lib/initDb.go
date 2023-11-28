package lib

import "database/sql"

func SetupDatabase(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user (
id INTEGER PRIMARY KEY AUTOINCREMENT,
username TEXT NOT NULL,
password TEXT NOT NULL,
deleted INTEGER NOT NULL DEFAULT 0
)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS session (
id INTEGER PRIMARY KEY AUTOINCREMENT,
token TEXT NOT NULL,
user_id INTEGER,
FOREIGN KEY(user_id) REFERENCES user(id) 
)`)

	if err != nil {
		return err
	}
	return nil
}
