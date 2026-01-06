package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/alejandrobyrne/website/views/home" // Import your views
)

func main() {
	// 1. Serve Static Files (CSS, Images)
	// This tells Go: "If URL starts with /static/, look in the static folder"
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Home Route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// We will create this component in the next step
		component := home.Index()
		templ.Handler(component).ServeHTTP(w, r)
	})

	log.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
