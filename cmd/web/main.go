package main

import (
	"log"
	"net/http"

	"github.com/gauravhegade/goshorty/internal/models/store"
	_ "github.com/mattn/go-sqlite3"
)

// dependency injection of URLDataModel
// making it available to all underlying handler functions
type app struct {
	store *store.URLDataModel
}

func main() {
	store, err := store.NewStore()
	if err != nil {
		log.Fatal(err)
		return
	}

	app := &app{
		store: store,
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Server running on port :8080")
	log.Fatal(server.ListenAndServe())
}
