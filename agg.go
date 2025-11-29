package main

import (
	"context"
	"fmt"
	"time"

	"github.com/vimto1234/gator/internal/database"
)

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		fmt.Print(err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	})

	if err != nil {
		fmt.Print(err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("Printing from %v\n", feed.Url)

	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
}
