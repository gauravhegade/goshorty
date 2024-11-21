package sqlite

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/gauravhegade/goshorty/internal/models"
)

const (
	shortURLlen = 6
)

type URLDataModel struct {
	DB *sql.DB
}

// returns all urls shortened from the database
func (m *URLDataModel) All() ([]models.URLData, error) {
	stmt := `
	SELECT short_code, long_url, title, created_on, expires_on
	FROM urls
	ORDER BY created_on DESC;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	urlData := []models.URLData{}
	for rows.Next() {
		var u models.URLData
		var expiresOn sql.NullTime
		if err := rows.Scan(&u.ShortCode, &u.URL, &u.Title, &u.CreatedOn, &expiresOn); err != nil {
			log.Fatal(err)
		}
		if expiresOn.Valid {
			u.ExpiresOn = &expiresOn.Time
		}
		urlData = append(urlData, u)
	}

	return urlData, nil
}

// creates a new short url
func (m *URLDataModel) CreateShortURL(c context.Context, url, title string, expiry time.Duration) (string, error) {
	shortCode := generateRandomString(shortURLlen)
	createdOn := time.Now()

	urlData := models.URLData{
		ShortCode: shortCode,
		URL:       url,
		Title:     title,
		CreatedOn: createdOn,
	}

	if expiry > 0 {
		expiresOn := createdOn.Add(expiry)
		urlData.ExpiresOn = &expiresOn
	}

	// TODO: DO NOT INSERT INTO DB FOR EACH AND EVERY SHORT URL
	// INSERT INTO DB IN BATCHES
	// TEMPORARILY CACHE IT INTO SOME CACHE STORE
	stmt := `
	INSERT INTO urls
	(short_code, long_url, title, created_on, expires_on)
	values (?, ?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt, urlData.ShortCode, urlData.URL, urlData.Title, urlData.CreatedOn, urlData.ExpiresOn)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}
