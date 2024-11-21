package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gauravhegade/goshorty/internal/models/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

// dependency injection of URLDataModel
// making it available to all underlying handler functions
type app struct {
	urls *sqlite.URLDataModel
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err.Error())
	}

	app := app{
		urls: &sqlite.URLDataModel{
			DB: db,
		},
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Server running on port :8080")
	log.Fatal(server.ListenAndServe())
}
