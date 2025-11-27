package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssFeed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rssFeed, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return &rssFeed, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &rssFeed, err
	}

	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return &rssFeed, err
	}

	rssFeed.Channel.Title = html.EscapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.EscapeString(rssFeed.Channel.Description)

	for i, item := range rssFeed.Channel.Item {
		item.Title = html.EscapeString(item.Title)
		item.Description = html.EscapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil

}
