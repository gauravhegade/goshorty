package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.getHome)
	mux.HandleFunc("GET /urls", app.fetchAllURLs)
	mux.HandleFunc("POST /shorten", app.shortenURL)
	mux.HandleFunc("GET /{shortCode}", app.redirectToURL)

	return mux
}
