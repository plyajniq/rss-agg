package main

import (
	"context"
	"log"
	"time"

	"github.com/plyajniq/rss-agg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequets time.Duration,
) {
	log.Printf("Start scrappig with %v goroutines and %s time between requests", concurrency, timeBetweenRequets)
	ticker := time.NewTicker(timeBetweenRequets)
	for ; ; <-ticker.C {
		feeds, err :=db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("error to fetching: %v", err)
			continue
		}
	}
}
