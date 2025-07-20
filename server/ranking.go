package server

import (
	"math"
	"strings"
	"time"
)

const (
	tfWeight       = 1.0
	filenameWeight = 2.0
	pathWeight     = 1.5
	recencyWeight  = 1.2
)

var filetypeWeight = map[string]float64{
	"md":  1.0,
	"txt": 1.0,
	"log": 0.2,
}

func RankFile(tf int, filename, query, path, ext string, mtime time.Time) float64 {
	score := 0.0

	score += termFrequencyScore(tf) * tfWeight
	score += filenameMatchScore(filename, query) * filenameWeight
	score += pathMatchScore(path, query) * pathWeight
	score += recencyScore(mtime) * recencyWeight
	score += filetypeScore(ext)

	return score
}

func termFrequencyScore(count int) float64 {
	if count == 0 {
		return 0
	}
	return 1 + math.Log(float64(count))
}

func filenameMatchScore(filename, query string) float64 {
	filename = strings.ToLower(filename)
	query = strings.ToLower(query)
	if filename == query {
		return 1
	}
	if strings.HasPrefix(filename, query) {
		return 0.9
	}
	if strings.Contains(filename, query) {
		return 0.7
	}
	return 0
}

func pathMatchScore(path, query string) float64 {
	path = strings.ToLower(path)
	query = strings.ToLower(query)
	if path == query {
		return 1
	}
	if strings.HasPrefix(path, query) {
		return 0.9
	}
	if strings.Contains(path, query) {
		return 0.7
	}
	return 0
}

func recencyScore(mtime time.Time) float64 {
	daysAgo := time.Since(mtime).Hours() / 24
	switch {
	case daysAgo <= 1:
		return 1
	case daysAgo <= 7:
		return 0.8
	case daysAgo <= 30:
		return 0.5
	default:
		return 0.2
	}
}

func filetypeScore(ext string) float64 {
	if score, ok := filetypeWeight[ext]; ok {
		return score
	}
	if len(ext) == 0 {
		return 0
	}
	return 0.8
}
