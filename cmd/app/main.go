package main

import (
	"fmt"
	"net/http"

	"github.com/captain-corgi/go-shorten-url/pkg/shortener"
)

func main() {
	urlSh := shortener.NewURLShortener()

	http.HandleFunc("/", urlSh.HandleForm)
	http.HandleFunc("/shorten", urlSh.HandleShorten)
	http.HandleFunc("/short/", urlSh.HandleRedirect)

	fmt.Println("URL Shortener is running on :8080")
	http.ListenAndServe(":8080", nil)
}
