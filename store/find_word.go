package store

import (
	"database/sql"
	"fmt"
	"time"
)

type FileWordInfo struct {
	Word      string
	Path      string
	Frequency int
	Size      int64
	Mtime     time.Time
	IsText    bool
	Extension string
}

func FindWordFileInfo(db *sql.DB, word string) ([]FileWordInfo, error) {
	q := `SELECT files.path, files.size, files.mtime, files.is_text, files.extension,
             CASE WHEN inverted_index.word LIKE ? THEN inverted_index.word ELSE NULL END,
             CASE WHEN inverted_index.word LIKE ? THEN inverted_index.freq ELSE NULL END
      FROM files
      LEFT JOIN inverted_index ON files.id = inverted_index.file_id
      WHERE inverted_index.word LIKE ? OR files.path LIKE ?;`
	wordLike := fmt.Sprintf("%%%s%%", word)
	rows, err := db.Query(q, wordLike, wordLike, wordLike, wordLike)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []FileWordInfo
	for rows.Next() {
		var info FileWordInfo
		var word sql.NullString
		var freq sql.NullInt64
		err := rows.Scan(&info.Path,
			&info.Size,
			&info.Mtime,
			&info.IsText,
			&info.Extension,
			&word,
			&freq,
		)
		if err != nil {
			return nil, err
		}
		if word.Valid {
			info.Word = word.String
		}
		if freq.Valid {
			info.Frequency = int(freq.Int64)
		}
		results = append(results, info)
	}
	return results, nil
}
