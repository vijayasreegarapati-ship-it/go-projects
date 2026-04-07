package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Link struct {
	ID        string `json:"id"`
	LongURL   string `json:"longUrl"`
	ShortCode string `json:"shortCode"`
	Clicks    int    `json:"clicks"`
}

var db = make(map[string]Link)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db["vue"] = Link{ID: "1", LongURL: "https://vuejs.org", ShortCode: "vue", Clicks: 42}

	http.HandleFunc("/api/links", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getLinks(w, r)
		case "POST":
			createLink(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/links/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			deleteLink(w, r)
		case "PUT":
			updateLink(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/", enableCORS(handleRedirect))

	fmt.Println("Full CRUD Go server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func getLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var linksArray []Link
	for _, link := range db {
		linksArray = append(linksArray, link)
	}
	json.NewEncoder(w).Encode(linksArray)
}

func createLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newLink Link
	json.NewDecoder(r.Body).Decode(&newLink)

	newLink.ShortCode = generateShortCode(5)
	newLink.ID = fmt.Sprintf("%d", time.Now().UnixNano()) // Simple unique ID
	newLink.Clicks = 0

	db[newLink.ShortCode] = newLink

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newLink)
}

func deleteLink(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/api/links/")

	delete(db, shortCode)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Link deleted"}`)
}

func updateLink(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/api/links/")

	existingLink, exists := db[shortCode]
	if !exists {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	var updatedData Link
	json.NewDecoder(r.Body).Decode(&updatedData)

	existingLink.LongURL = updatedData.LongURL
	db[shortCode] = existingLink

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingLink)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	link, exists := db[shortCode]
	if !exists {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	link.Clicks++
	db[shortCode] = link

	http.Redirect(w, r, link.LongURL, http.StatusFound)
}

func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
