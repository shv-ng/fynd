package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shv-ng/fynd/app"
	"github.com/shv-ng/fynd/crawler"
	"github.com/shv-ng/fynd/indexer"
	"github.com/shv-ng/fynd/store"
	"github.com/shv-ng/fynd/types"
)

func Sync(ctx app.Context) {
	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, ctx.Setting.MaxConcurrency)
	crawlerToIndexerCh := make(chan types.File, 1000)
	indexerToDBCh := make(chan types.IndexedFile, 1000)

	dbcache, err := store.DBCache(ctx.DB)
	if err != nil {
		log.Fatal(err)
	}

	crl := crawler.Crawler{
		DBCache:  dbcache,
		Settings: ctx.Setting,
		Mu:       &mu,
		Wg:       &wg,
		Sem:      sem,
		Ch:       crawlerToIndexerCh,
	}
	wg.Add(1)
	go crl.Crawl(ctx.Setting.RootPath)
	go func() {
		wg.Wait()
		close(crawlerToIndexerCh)
	}()
	countUpdated, deletedPath := indexer.Indexer(crl.DBCache, crawlerToIndexerCh, indexerToDBCh)

	if err := store.RemoveDeletedFiles(deletedPath, ctx.DB); err != nil {
		log.Fatalln(err)
	}
	if err = store.BatchInsertHandler(indexerToDBCh, ctx.DB); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("📁 Crawled: %v files ✏️ Updated/Inserted: %v files ⏱️ Duration: %v\n\n", crl.CountCrawled, countUpdated, time.Since(start))
}
