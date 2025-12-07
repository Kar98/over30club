package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	root := flag.String("root", ".", "project root (where artistdata and display live)")
	flag.Parse()

	artistDir := filepath.Join(*root, "artistdata")
	displayDir := filepath.Join(*root, "display")

	// Serve artist files (so the browser can fetch individual artist JSONs)
	http.Handle("/artistdata/", addCORS(http.StripPrefix("/artistdata/", http.FileServer(http.Dir(artistDir)))))

	// Serve the display directory (your HTML + JS/CSS)
	http.Handle("/display/", addCORS(http.StripPrefix("/display/", http.FileServer(http.Dir(displayDir)))))

	// API endpoint to list artist files
	http.HandleFunc("/api/artist-files", func(w http.ResponseWriter, r *http.Request) {
		// handle preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// set CORS for actual response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		artistFilesHandler(w, r, artistDir)
	})

	addr := ":" + *port
	log.Printf("Serving on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func artistFilesHandler(w http.ResponseWriter, r *http.Request, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files := []string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		files = append(files, name)
	}
	w.Header().Set("Content-Type", "application/json")
	// also set CORS header just in case
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(files)
}

// addCORS wraps a handler and sets permissive CORS headers.
func addCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
