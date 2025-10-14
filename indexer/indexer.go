package indexer

import (
	"time"

	"github.com/shv-ng/fynd/models"
)

func Indexer(dbcache map[string]time.Time, input <-chan models.File, output chan<- models.IndexedFile) (int, []string) {
	count := 0
	for file := range input {
		count++
		var invertedIndexes []models.InvertedIndex
		var wordCount map[string]int
		if file.IsText {
			wordCount = map[string]int{}
			content := file.Content
			words := Sanatize(content)
			for _, w := range words {
				wordCount[w]++
			}
		} else {
			wordCount = nil
		}
		for word, freq := range wordCount {
			invertedIndexes = append(invertedIndexes, models.InvertedIndex{
				Word:      word,
				FileID:    0,
				Frequency: freq,
			})
		}
		indexed := models.IndexedFile{
			File:            file,
			InvertedIndexes: invertedIndexes,
		}
		output <- indexed
	}

	close(output)
	// return deleted files
	var paths []string
	for path := range dbcache {
		paths = append(paths, path)
	}
	return count, paths
}
