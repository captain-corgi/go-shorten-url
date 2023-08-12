package shortener

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type (
	URLShortener interface {
		HandleForm(w http.ResponseWriter, rq *http.Request)
		HandleShorten(w http.ResponseWriter, rq *http.Request)
		HandleRedirect(w http.ResponseWriter, rq *http.Request)
	}
	URLShortenerImpl struct {
		urls map[string]string
	}
)

func NewURLShortener() *URLShortenerImpl {
	return &URLShortenerImpl{
		urls: make(map[string]string),
	}
}

func (r *URLShortenerImpl) HandleForm(w http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodPost {
		http.Redirect(w, rq, "/shorten", http.StatusSeeOther)
		return
	}

	// Serve the HTML form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
	`)
}

func (r *URLShortenerImpl) HandleShorten(w http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	originalURL := rq.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Generate a unique shortened key for the original URL
	shortKey := generateShortKey()
	r.urls[shortKey] = originalURL

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	// Render the HTML response with the shortened URL
	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
		<h2>URL Shortener</h2>
		<p>Original URL: %s</p>
		<p>Shortened URL: <a href="%s">%s</a></p>
		<form method="post" action="/shorten">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="submit" value="Shorten">
		</form>
	`, originalURL, shortenedURL, shortenedURL)
	fmt.Fprintf(w, responseHTML)
}

func (r *URLShortenerImpl) HandleRedirect(w http.ResponseWriter, rq *http.Request) {
	shortKey := rq.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := r.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	// Redirect the user to the original URL
	http.Redirect(w, rq, originalURL, http.StatusMovedPermanently)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	return string(shortKey)
}
