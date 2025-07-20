package store

import (
	"database/sql"

	_ "github.com/glebarez/sqlite"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return nil, err
	}

	// Pings to verify connection (important!)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	schema := `
    CREATE TABLE IF NOT EXISTS files (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        path TEXT UNIQUE,
        size INTEGER,
        mtime DATETIME,
        is_text BOOLEAN,
        extension TEXT
    );
    CREATE TABLE IF NOT EXISTS inverted_index (
        word TEXT,
        file_id INTEGER,
				freq INTEGER,
        FOREIGN KEY(file_id) REFERENCES files(id) ON DELETE CASCADE
    );`
	_, err = db.Exec(schema)

	return db, err
}
