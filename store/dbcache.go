package store

import (
	"database/sql"
	"time"
)

func DBCache(db *sql.DB) (map[string]time.Time, error) {
	// file path : mtime
	cache := make(map[string]time.Time)
	q := "SELECT path,mtime FROM files;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			path  string
			mtime time.Time
		)
		rows.Scan(&path, &mtime)
		cache[path] = mtime
	}
	return cache, nil
}
