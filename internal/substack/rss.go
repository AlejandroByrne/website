package substack

import (
	"github.com/mmcdole/gofeed"
)

type Post struct {
	Title string
	Link  string
	Date  string
}

// FetchFeed gets the top N posts from a substack URL
func FetchFeed(url string, limit int) ([]Post, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for i, item := range feed.Items {
		if i >= limit {
			break
		}
		// Default to empty string if Published is missing
		date := item.Published
		if date == "" {
			date = "Recently"
		}

		posts = append(posts, Post{
			Title: item.Title,
			Link:  item.Link,
			Date:  date,
		})
	}
	return posts, nil
}

// FetchFeedCached wraps FetchFeed with an in-memory TTL cache.
func FetchFeedCached(c *FeedCache, url string, limit int) ([]Post, error) {
	if c != nil {
		if posts, ok := c.Get(url, limit); ok {
			return posts, nil
		}
	}
	posts, err := FetchFeed(url, limit)
	if err != nil {
		return nil, err
	}
	if c != nil {
		c.Set(url, limit, posts)
	}
	return posts, nil
}
