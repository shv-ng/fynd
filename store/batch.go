package store

import (
	"database/sql"

	"github.com/shv-ng/fynd/types"
)

const batchSize = 100

func BatchInsertHandler(ch <-chan types.IndexedFile, db *sql.DB) error {
	var files []types.IndexedFile
	for file := range ch {
		files = append(files, file)

		if len(files) >= batchSize {

			if err := batchInsert(files, db); err != nil {
				return err
			}

			files = nil
		}
	}
	// insert leftover files
	if len(files) > 0 {
		if err := batchInsert(files, db); err != nil {
			return err
		}
	}
	return nil
}

func batchInsert(files []types.IndexedFile, db *sql.DB) error {
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

	fileStmt, err := tx.Prepare(`
	INSERT INTO files (path, size, mtime, is_text, extension)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(path) DO UPDATE SET
		size = excluded.size,
		mtime = excluded.mtime,
		is_text = excluded.is_text,
		extension = excluded.extension
`)
	if err != nil {
		return err
	}

	defer fileStmt.Close()

	indexStmt, err := tx.Prepare(`
		INSERT INTO inverted_index (word, file_id, freq) 
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer indexStmt.Close()

	for _, f := range files {
		// Insert or update file
		_, err := fileStmt.Exec(f.File.Path, f.File.Size, f.File.MTime, f.File.IsText, f.File.Extension)
		if err != nil {
			return err
		}

		// Get the file ID
		var fileID int
		err = tx.QueryRow(`SELECT id FROM files WHERE path = ?`, f.File.Path).Scan(&fileID)
		if err != nil {
			return err
		}

		// Delete old inverted indexes for this file
		_, err = tx.Exec(`DELETE FROM inverted_index WHERE file_id = ?`, fileID)
		if err != nil {
			return err
		}

		// Insert new inverted indexes
		for _, idx := range f.InvertedIndexes {
			_, err := indexStmt.Exec(idx.Word, fileID, idx.Frequency)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
