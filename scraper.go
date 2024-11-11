package main

import (
	"context"
	"log"
	"sync"
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
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("error to fetching: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		// inter feeds in main goroutine
		for _, feed := range feeds {
			// +1 to counter of goroutines to wait
			wg.Add(1)
			// spawn goroutine
			go scrapeFeed(db, wg, feed)
		}
		// wait for all to finish
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	// -1 to counter of goroutines to wait
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error to fetching: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Fail to get RSSFeed from URL: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post: ", item.Title, "from ", feed.Name)
	}
	log.Printf("Feed %s updated with %v posts", feed.Name, len(rssFeed.Channel.Item))

}
