package projects_store

import "strings"

type Project struct {
	ID          string
	Title       string
	Description string
	TechStack   []string
	Link        string
}

// Your Mock Database
var allProjects = []Project{
	{
		ID:          "1",
		Title:       "Azure RAG Chatbot",
		Description: "An AI chatbot that uses Retrieval-Augmented Generation to query SharePoint documents.",
		TechStack:   []string{"Azure OpenAI", "Python", "React", "CosmosDB"},
		Link:        "https://github.com/alejandrobyrne/rag-bot",
	},
	{
		ID:          "2",
		Title:       "Distributed Task Queue",
		Description: "A high-performance task processing pipeline handling 10k events/sec.",
		TechStack:   []string{"Go", "RabbitMQ", "Docker", "gRPC"},
		Link:        "https://github.com/alejandrobyrne/task-queue",
	},
	{
		ID:          "3",
		Title:       "Renaissance Art Gallery",
		Description: "A digital gallery showcasing classical art, utilizing high-res image tiling.",
		TechStack:   []string{"TypeScript", "WebGL", "Next.js"},
		Link:        "#",
	},
}

func Search(query string) []Project {
	if query == "" {
		return allProjects
	}

	var results []Project
	query = strings.ToLower(query)

	for _, p := range allProjects {
		if strings.Contains(strings.ToLower(p.Title), query) ||
			strings.Contains(strings.ToLower(p.Description), query) {
			results = append(results, p)
		}
	}
	return results
}
