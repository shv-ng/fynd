package models

import "time"

type File struct {
	Path      string
	Content   string
	Size      int64
	MTime     time.Time
	IsText    bool
	Extension string
}
type IndexedFile struct {
	File            File
	InvertedIndexes []InvertedIndex
}
