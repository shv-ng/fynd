package server

import (
	"fmt"
	"log"
	"math"
	"path/filepath"
	"slices"
	"sort"
	"time"

	"github.com/ShivangSrivastava/fynd/app"
	"github.com/ShivangSrivastava/fynd/store"
)

// fileInfo aggregates file metadata and word frequency data.
type fileInfo struct {
	Size   int64
	Mtime  time.Time
	IsText bool
	Ext    string
	Words  map[string]int
}

// wordFreq pairs a word with its frequency.
type wordFreq struct {
	word string
	freq int
}

// Find performs a file search based on a parsed query and outputs ranked results.
func Find(ctx app.Context, input string, opts QueryOptions) {
	opts = ParseQuery(ctx, input, opts)

	fileMap := make(map[string]*fileInfo)
	rankingMap := make(map[string]float64)
	var filePaths []string

	for _, word := range opts.Query {
		infos, err := store.FindWordFileInfo(ctx.DB, word)
		if err != nil {
			log.Printf("FindWordFileInfo error: %v\n", err)
			continue
		}

		for _, info := range infos {
			if len(opts.Ext) > 0 && !slices.Contains(opts.Ext, info.Extension) {
				continue
			}

			if _, exists := fileMap[info.Path]; !exists {
				fileMap[info.Path] = &fileInfo{
					Size:   info.Size,
					Mtime:  info.Mtime,
					IsText: info.IsText,
					Ext:    info.Extension,
					Words:  make(map[string]int),
				}
				filePaths = append(filePaths, info.Path)
			}

			file := fileMap[info.Path]
			file.Words[info.Word] = info.Frequency
			rankingMap[info.Path] += RankFile(info.Frequency, filepath.Base(info.Path), info.Word, filepath.Dir(info.Path), info.Extension, info.Mtime)
		}
	}

	sort.Slice(filePaths, func(i, j int) bool {
		return rankingMap[filePaths[i]] > rankingMap[filePaths[j]]
	})

	if opts.Top > 0 && len(filePaths) > opts.Top {
		filePaths = filePaths[:opts.Top]
	}

	printResults(filePaths, fileMap, rankingMap)
}

// printResults displays ranked file results with metadata and word frequency.
func printResults(filePaths []string, fileMap map[string]*fileInfo, rankingMap map[string]float64) {
	for i, path := range filePaths {
		info := fileMap[path]
		fmt.Printf("\033[1;34m[%d]\033[0m \033[1;36m%s\033[0m\n", i+1, path)
		fmt.Printf("    \033[1;33msize:\033[0m %dB  \033[1;33mmodified:\033[0m %s  \033[1;33mtype:\033[0m %s(%s)  \033[1;33mscore:\033[0m %.4f\n",
			info.Size,
			info.Mtime.Format("2006-01-02 15:04:05"),
			colorizeType(info.IsText),
			info.Ext,
			roundScore(rankingMap[path]),
		)

		sortedWords := sortWordsByFrequency(info.Words)
		fmt.Printf("    \033[1;33mwords:\033[0m")
		for _, wf := range sortedWords {
			if wf.freq > 0 {
				fmt.Printf(" \033[1;37m%s\033[0m:\033[90m%d\033[0m", wf.word, wf.freq)
			}
		}
		fmt.Println()
	}
}

// sortWordsByFrequency returns a sorted slice of wordFreq by descending frequency.
func sortWordsByFrequency(words map[string]int) []wordFreq {
	var sorted []wordFreq
	for w, f := range words {
		sorted = append(sorted, wordFreq{w, f})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].freq > sorted[j].freq
	})
	return sorted
}

// colorizeType returns a color-coded string representing file type.
func colorizeType(isText bool) string {
	if isText {
		return "\033[32mtext\033[0m"
	}
	return "\033[31mbinary\033[0m"
}

// roundScore rounds a float64 score to four decimal places.
func roundScore(score float64) float64 {
	return math.Round(score*10000) / 10000
}
