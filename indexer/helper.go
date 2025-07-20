package indexer

import (
	"regexp"
	"strings"

	"github.com/dchest/stemmer/porter2"
)

func Sanatize(content string) []string {
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
