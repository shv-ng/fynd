package store

import "database/sql"

func RemoveDeletedFiles(filePaths []string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt, err := tx.Prepare(`DELETE FROM files WHERE path = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, path := range filePaths {
		_, err := stmt.Exec(path)
		if err != nil {
			return err
		}
	}

	return nil
}
