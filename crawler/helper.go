package crawler

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/ShivangSrivastava/fynd/models"
)

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
func processFile(path string, info os.FileInfo) (models.File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return models.File{}, err
	}

	isText := utf8.Valid(data)
	content := ""
	if isText {
		content = string(data)
	}

	ext := strings.TrimPrefix(filepath.Ext(info.Name()), ".")
	return models.File{
		Path:      path,
		Content:   content,
		Size:      info.Size(),
		MTime:     info.ModTime(),
		IsText:    isText,
		Extension: ext,
	}, nil
}
