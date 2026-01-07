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
