package crawler

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/shv-ng/fynd/app"
	"github.com/shv-ng/fynd/types"
)

type Crawler struct {
	DBCache      map[string]time.Time
	Settings     app.Settings
	CountCrawled int
	Mu           *sync.Mutex
	Wg           *sync.WaitGroup
	Sem          chan struct{}
	Ch           chan<- types.File
}

// Crawl in given path, respect the config, choose file accordingly and pass to
// channel
func (c *Crawler) Crawl(path string) {
	defer c.Wg.Done()
	// list out all entries from given dir path
	entries, err := c.readDirSafe(path)
	if err != nil {
		log.Printf("error reading dir %s: %v", path, err)
		return
	}
	// traverse on path
	for _, entry := range entries {

		name := entry.Name()
		absPath := filepath.Join(path, name)
		info, err := entry.Info()
		if err != nil {
			log.Printf("error on getting info for %s: %v", absPath, err)
			continue
		}
		// found dir, crawl again
		if info.IsDir() {
			if c.shouldCrawlDir(name) {
				c.Wg.Add(1)
				go c.Crawl(absPath)
			}
			continue
		}
		// remove seen paths from dbcache
		c.Mu.Lock()
		mtime, ok := c.DBCache[absPath]
		if ok {
			delete(c.DBCache, absPath)
		}
		c.CountCrawled++
		c.Mu.Unlock()
		// found same file with no modification, go to next entry
		if mtime.Equal(info.ModTime()) {
			continue
		}
		// read all the info, contents from file
		file, err := processFile(absPath, info)
		if err != nil {
			log.Printf("error on processing file %s: %v", absPath, err)
			continue
		}
		c.Ch <- file
	}
}

// read dir on keeping concurrency limit
func (c *Crawler) readDirSafe(path string) ([]os.DirEntry, error) {
	c.Sem <- struct{}{}
	defer func() { <-c.Sem }()
	return os.ReadDir(path)
}

// check isnt it excluded or any other path or not
func (c *Crawler) shouldCrawlDir(name string) bool {
	if !c.Settings.IncludeHidden && strings.HasPrefix(name, ".") {
		return false
	}
	if slices.Contains(c.Settings.ExcludeDirs, name) {
		return false
	}
	return len(c.Settings.IncludeDirs) == 0 || slices.Contains(c.Settings.IncludeDirs, name)
}

// get file info and read contents
func processFile(path string, info os.FileInfo) (types.File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.File{}, err
	}

	isText := utf8.Valid(data)
	content := ""
	if isText {
		content = string(data)
	}

	ext := strings.TrimPrefix(filepath.Ext(info.Name()), ".")
	return types.File{
		Path:      path,
		Content:   content,
		Size:      info.Size(),
		MTime:     info.ModTime(),
		IsText:    isText,
		Extension: ext,
	}, nil
}
