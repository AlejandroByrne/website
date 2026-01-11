package projects_store

import (
	"encoding/csv"
	"net/http"
	"strings"
)

const ProjectsSheetURL = "https://docs.google.com/spreadsheets/d/e/2PACX-1vRW6dpD7e9WbKuBeN-b8w2SgkPI_uLucHFavnd0koZrBBtvTqysPfGzANt4kNHPL5bSgtOcqSDw_-cY/pub?output=csv"

type Project struct {
	ID          string
	Title       string
	Description string
	TechStack   []string
	Link        string
}

// FetchProjects retrieves projects from the Google Sheet
func FetchProjects(query string) ([]Project, error) {
	resp, err := http.Get(ProjectsSheetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) < 1 {
		return []Project{}, nil
	}
	dataRows := records[1:]

	var results []Project
	query = strings.ToLower(query)

	for i, row := range dataRows {
		// Expecting columns: Title, Desc, Topics, Url
		if len(row) < 4 {
			continue
		}

		title := row[0]
		desc := row[1]
		topicsStr := row[2]
		link := row[3]

		// Filter if query is present
		if query != "" {
			if !strings.Contains(strings.ToLower(title), query) &&
				!strings.Contains(strings.ToLower(desc), query) &&
				!strings.Contains(strings.ToLower(topicsStr), query) {
				continue
			}
		}

		// Parse topics
		var techStack []string
		if topicsStr != "" {
			rawTopics := strings.Split(topicsStr, ",")
			for _, t := range rawTopics {
				techStack = append(techStack, strings.TrimSpace(t))
			}
		}

		// Generate a simple ID based on index
		id := string(rune('1' + i)) // '1', '2', etc. - simple for now

		results = append(results, Project{
			ID:          id,
			Title:       title,
			Description: desc,
			TechStack:   techStack,
			Link:        link,
		})
	}

	// Optional: Reverse to show newest first if the sheet is ordered chronologically?
	// Or keep sheet order. Let's keep sheet order for manual control.
	return results, nil
}

// Search is now a wrapper around FetchProjects for compatibility or specific filtering logic
// But since we are fetching live, we just call FetchProjects with the query.
func Search(query string) []Project {
	projects, err := FetchProjects(query)
	if err != nil {
		// Log error in a real app, for now return empty
		return []Project{}
	}
	return projects
}
