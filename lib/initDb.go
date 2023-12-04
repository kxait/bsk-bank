package lib

import "database/sql"

func SetupDatabase(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS user (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT NOT NULL,
      password TEXT NOT NULL,
      deleted INTEGER NOT NULL DEFAULT 0
)`,
		`CREATE TABLE IF NOT EXISTS session (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      token TEXT NOT NULL,
      user_id INTEGER,
      created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
      expires_at TEXT NOT NULL,
      valid INTEGER NOT NULL DEFAULT 1,
      FOREIGN KEY(user_id) REFERENCES user(id) 
)`,
		`CREATE TABLE IF NOT EXISTS account (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      holder_name TEXT NOT NULL,
      created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
      deleted INTEGER NOT NULL DEFAULT 0
)`,
		`CREATE TABLE IF NOT EXISTS account_transaction (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      source_account_id INTEGER NOT NULL,
      destination_account_id INTEGER NOT NULL,
      amount REAL NOT NULL,
      FOREIGN KEY (source_account_id) REFERENCES account(id),
      FOREIGN KEY (destination_account_id) REFERENCES account(id)
)`,
		`CREATE TABLE IF NOT EXISTS failed_login (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT NOT NULL,
      created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
      ip_address TEXT NOT NULL
)`,
		`CREATE TABLE IF NOT EXISTS config (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      key TEXT NOT NULL,
      value TEXT NOT NULL,
      modified_at TEXT
)`,
		`INSERT OR IGNORE INTO account (id, holder_name) VALUES (-1, 'Okienko')`,
		`INSERT OR IGNORE INTO config (id, key, value) VALUES (0, 'transaction_limit', '100000')`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)

		if err != nil {
			return err
		}
	}
	return nil
}
