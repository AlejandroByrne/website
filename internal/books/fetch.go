package books

import (
	"encoding/csv"
	"net/http"
)

const SheetURL = "https://docs.google.com/spreadsheets/d/e/2PACX-1vRJPJIuZ3d61Kzxr8zEkyyKvcjyH0t-kXMBp1bFjht5rOlcumx-DODqGgRYnRXbaIjyIqzUOriN86rm/pub?output=csv"

type Book struct {
	ISBN       string
	Title      string
	Author     string
	Translator string
	Pages      string // Keeping as string is fine for display
	Year       string
	Topics     string
	CoverURL   string
	Thoughts   string
	ReadDate   string
}

func FetchRecent(limit int) ([]Book, error) {
	resp, err := http.Get(SheetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip header row
	if len(records) < 1 {
		return []Book{}, nil
	}
	dataRows := records[1:]

	// Parse into structs
	var books []Book
	for _, row := range dataRows {
		// It expects 10 rows per book column (now including ReadDate)
		if len(row) < 10 {
			continue
		}
		b := Book{
			ISBN:       row[0],
			Title:      row[1],
			Author:     row[2],
			Translator: row[3],
			Pages:      row[4],
			Year:       row[5],
			Topics:     row[6],
			CoverURL:   row[7],
			Thoughts:   row[8],
			ReadDate:   row[9],
		}
		books = append(books, b)
	}
	// Reverse list to get newest first
	for i, j := 0, len(books)-1; i < j; i, j = i+1, j-1 {
		books[i], books[j] = books[j], books[i]
	}
	// Apply limit
	if len(books) > limit {
		books = books[:limit]
	}

	return books, nil
}
