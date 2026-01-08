package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"

	// IMPORTS
	"github.com/alejandrobyrne/website/internal/projects_store"
	"github.com/alejandrobyrne/website/internal/substack"
	"github.com/alejandrobyrne/website/views/about"
	"github.com/alejandrobyrne/website/views/home"
	"github.com/alejandrobyrne/website/views/projects_view"
)

func main() {
	// Serve Static Files (CSS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// --- ROUTE 1: HOME (DASHBOARD) ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// 1. Fetch Substack (Top 3)
		posts, err := substack.FetchFeed("https://rubylarocca.substack.com/feed", 3)
		if err != nil {
			// Fail gracefully
			posts = []substack.Post{}
			log.Println("Error fetching feed:", err)
		}

		// 2. Fetch Projects (Top 2)
		allProjects := projects_store.Search("")
		limit := 2
		if len(allProjects) < limit {
			limit = len(allProjects)
		}
		featuredProjects := allProjects[:limit]

		// 3. Render
		data := home.HomeData{
			RecentPosts:      posts,
			FeaturedProjects: featuredProjects,
		}
		component := home.Index(data)
		templ.Handler(component).ServeHTTP(w, r)
	})

	// --- ROUTE 2: FULL PROJECTS PAGE ---
	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		data := projects_store.Search("")
		component := projects_view.Page(data)
		templ.Handler(component).ServeHTTP(w, r)
	})

	// --- ROUTE 3: PROJECTS SEARCH (HTMX) ---
	http.HandleFunc("/projects/search", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		query := r.FormValue("query")

		filtered := projects_store.Search(query)

		component := projects_view.ProjectList(filtered)
		templ.Handler(component).ServeHTTP(w, r)
	})

	// --- ROUTE 4: ABOUT PAGE ---
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		component := about.Index()
		templ.Handler(component).ServeHTTP(w, r)
	})

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
