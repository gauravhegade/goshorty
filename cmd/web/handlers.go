package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gauravhegade/goshorty/internal/models/store"
)

type httpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type shortURLRequest struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
	// nil value exists for int64
	// but I want to completely ignore the field if not provided
	ExpiryInSecs *int64 `json:"expiry_in_secs"`
}

func (app *app) sendResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(httpResponse{Status: "success", Data: data})
	if err != nil {
		app.sendErrorResponse(w, "Internal server error!", http.StatusInternalServerError, nil)
		return
	}
	w.Write(b)
}

func (app *app) sendErrorResponse(w http.ResponseWriter, message string, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	b, _ := json.Marshal(httpResponse{Status: "error", Message: message, Data: data})
	w.Write(b)
}

func (app *app) getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	app.sendResponse(w, "Welcome to homepage")
}

func (app *app) fetchAllURLs(w http.ResponseWriter, r *http.Request) {
	urls, err := app.store.All()
	if err != nil {
		log.Fatal(err)
	}

	app.sendResponse(w, map[string]any{
		"URLs": urls,
	})
}

func (app *app) shortenURL(w http.ResponseWriter, r *http.Request) {
	var shortURLReq shortURLRequest

	if err := json.NewDecoder(r.Body).Decode(&shortURLReq); err != nil {
		app.sendErrorResponse(w, "Invalid request body!", http.StatusBadRequest, nil)
		return
	}

	// url validation
	if valid := app.isValidURL(shortURLReq.URL); !valid {
		app.sendErrorResponse(w, "Invalid URL provided!", http.StatusBadRequest, nil)
		return
	}

	// log.Println(shortURLReq.ExpiryInSecs)
	var expiry time.Duration
	if shortURLReq.ExpiryInSecs != nil && *shortURLReq.ExpiryInSecs > 0 {
		expiry = time.Duration(*shortURLReq.ExpiryInSecs) * time.Second
	}

	// log.Println(expiry)
	shortCode, err := app.store.CreateShortURL(shortURLReq.URL, shortURLReq.Title, expiry)
	if err != nil {
		app.sendErrorResponse(w, "Internal server error!"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	app.sendResponse(w, map[string]any{
		"ShortCode": shortCode,
	})
}

func (app *app) redirectToURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	if shortCode == "" {
		app.sendErrorResponse(w, "Invalid short code!", http.StatusBadRequest, nil)
	}

	urlData, err := app.store.GetRedirectData(shortCode)
	if err != nil {
		if err == store.ErrNotExists {
			log.Println(err.Error())
			app.sendErrorResponse(w, "URL not found", http.StatusNotFound, nil)
			return
		} else if err == store.ErrURLExpired {
			log.Println(err.Error())
			app.sendErrorResponse(w, "URL expired", http.StatusGone, nil)
			return
		}
	}

	// check if url has some prefix or not
	// if not then add some prefix
	if !strings.HasPrefix(urlData.URL, "https://") && !strings.HasPrefix(urlData.URL, "http://") {
		urlData.URL = "https://" + urlData.URL
	}
	log.Println("Redirecting to", urlData.URL)
	w.Header().Set("Location", urlData.URL)
	w.WriteHeader(http.StatusFound)
}

func (app *app) deleteURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	if shortCode == "" {
		app.sendErrorResponse(w, "Invalid short code!", http.StatusBadRequest, nil)
	}

	err := app.store.DeleteURLData(shortCode)
	if err != nil {
		if err == store.ErrNotExists {
			log.Println(err.Error())
			app.sendErrorResponse(w, "URL not found", http.StatusNotFound, nil)
			return
		}
		log.Println(err.Error())
		app.sendErrorResponse(w, "Internal server error!", http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.sendResponse(w, map[string]any{
		"message": "deleted requested resource",
	})
}

func (app *app) isValidURL(url string) bool {
	regex, _ := regexp.Compile(`[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	return regex.MatchString(url)
}
