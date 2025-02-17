package main

import (
	"log"
	"net/http"
	"os"
)

const (
	defaultPort = "3000"
)

func main() {
	http.HandleFunc("/pdf", wrapper(pdfHandler))
	http.HandleFunc("/docx", wrapper(docxHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Start listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func wrapper(h func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Use POST method", http.StatusMethodNotAllowed)
			return
		}
		if err := h(w, r); err != nil {
			log.Printf("Request failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
