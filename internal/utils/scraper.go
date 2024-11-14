package utils

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"rss-agg/internal/database"

	"github.com/google/uuid"
)

func StartScraping(
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

	rssFeed, err := ParseUrlToFeed(feed.Url)
	if err != nil {
		log.Printf("Fail to get RSSFeed from URL: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Fail to parse date %v with error: %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Fail to create post: %v", err)
		}
	}
	log.Printf("Feed %s updated with %v posts", feed.Name, len(rssFeed.Channel.Item))

}
