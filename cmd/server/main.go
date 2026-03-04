package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"

	// IMPORTS
	"github.com/alejandrobyrne/website/internal/books"
	"github.com/alejandrobyrne/website/internal/projects_store"
	"github.com/alejandrobyrne/website/internal/substack"
	"github.com/alejandrobyrne/website/views/about"
	"github.com/alejandrobyrne/website/views/home"
	"github.com/alejandrobyrne/website/views/projects_view"
	"github.com/alejandrobyrne/website/views/now"
	"github.com/alejandrobyrne/website/views/uses"
)

func main() {
	// Config
	feedURL := getenv("SUBSTACK_FEED_URL", "https://alejandrobyrne.substack.com/feed")
	revalidateToken := os.Getenv("REVALIDATE_TOKEN") // optional; if empty, endpoint is disabled
	cacheTTL := getenvDuration("FEED_CACHE_TTL", 15*time.Minute)

	// Cache
	feedCache := substack.NewFeedCache(cacheTTL)

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Home
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		posts, err := substack.FetchFeedCached(feedCache, feedURL, 3)
		if err != nil {
			log.Println("Error fetching feed:", err)
			posts = []substack.Post{}
		}

		allProjects := projects_store.Search("")
		limit := 2
		if len(allProjects) < limit { limit = len(allProjects) }
		featuredProjects := allProjects[:limit]

		recentBooks, err := books.FetchRecent(3)
		if err != nil {
			log.Println("Error fetching books:", err)
			recentBooks = []books.Book{}
		}

		data := home.HomeData{RecentPosts: posts, FeaturedProjects: featuredProjects, RecentBooks: recentBooks}
		component := home.Index(data)
		templ.Handler(component).ServeHTTP(w, r)
	})

	// Revalidate endpoint: clears feed cache
	if revalidateToken != "" {
		http.HandleFunc("/admin/revalidate", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
			if r.URL.Query().Get("key") != revalidateToken { http.Error(w, "unauthorized", http.StatusUnauthorized); return }
			feedCache.InvalidateAll()
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
	}

	// Projects page
	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		data := projects_store.Search("")
		component := projects_view.Page(data)
		templ.Handler(component).ServeHTTP(w, r)
	})

	// Projects search (HTMX)
	http.HandleFunc("/projects/search", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		query := r.FormValue("query")
		filtered := projects_store.Search(query)
		component := projects_view.ProjectList(filtered)
		templ.Handler(component).ServeHTTP(w, r)
	})

	// About
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		component := about.Index()
		templ.Handler(component).ServeHTTP(w, r)
	})

	// Now
	http.HandleFunc("/now", func(w http.ResponseWriter, r *http.Request) {
		component := now.Page()
		templ.Handler(component).ServeHTTP(w, r)
	})

	// Uses
	http.HandleFunc("/uses", func(w http.ResponseWriter, r *http.Request) {
		component := uses.Page()
		templ.Handler(component).ServeHTTP(w, r)
	})

	// Redirect to substack
	http.HandleFunc("/substack", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://alejandrobyrne.substack.com", http.StatusMovedPermanently)
	})

	// robots.txt
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("User-agent: *\nAllow: /\nSitemap: https://alejandrobyrne.com/sitemap.xml\n"))
	})

	// sitemap.xml (static routes for now)
	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>https://alejandrobyrne.com/</loc></url>
  <url><loc>https://alejandrobyrne.com/projects</loc></url>
  <url><loc>https://alejandrobyrne.com/about</loc></url>
</urlset>`
		_, _ = w.Write([]byte(sitemap))
	})

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil { return d }
	}
	return def
}
