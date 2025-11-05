package indexer

import (
	"regexp"
	"strings"
	"time"

	"github.com/dchest/stemmer/porter2"
	"github.com/shv-ng/fynd/types"
)

func Indexer(dbcache map[string]time.Time, input <-chan types.File, output chan<- types.IndexedFile) (int, []string) {
	count := 0
	for file := range input {
		count++
		var invertedIndexes []types.InvertedIndex
		var wordCount map[string]int
		if file.IsText {
			wordCount = map[string]int{}
			content := file.Content
			words := sanatize(content)
			for _, w := range words {
				wordCount[w]++
			}
		} else {
			wordCount = nil
		}
		for word, freq := range wordCount {
			invertedIndexes = append(invertedIndexes, types.InvertedIndex{
				Word:      word,
				FileID:    0,
				Frequency: freq,
			})
		}
		indexed := types.IndexedFile{
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

func sanatize(content string) []string {
	return stemmer(removeStopWords(tokenise(content)))
}

func tokenise(content string) []string {
	re := regexp.MustCompile("[a-zA-Z0-9-_]{3,}")
	matches := re.FindAllString(content, -1)

	for i, word := range matches {
		matches[i] = strings.ToLower(word)
	}

	return matches
}

func removeStopWords(content []string) []string {
	var res []string
	for _, word := range content {
		if _, ok := stopwords[word]; !ok {
			res = append(res, word)
		}
	}
	return res
}

func stemmer(content []string) []string {
	eng := porter2.Stemmer
	var res []string
	for _, word := range content {
		res = append(res, eng.Stem(word))
	}
	return res
}
