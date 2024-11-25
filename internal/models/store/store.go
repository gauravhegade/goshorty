package store

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/gauravhegade/goshorty/internal/models"
)

type URLDataModel struct {
	db          *sql.DB
	cache       map[string]models.URLData
	shortURLLen int
}

var ErrNotExists = errors.New("URL does not exist")
var ErrURLExpired = errors.New("URL has been expired")

func NewStore() (*URLDataModel, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return nil, err
	}

	urlDataModel := &URLDataModel{
		db:          db,
		cache:       make(map[string]models.URLData),
		shortURLLen: 6,
	}

	// load the cache, helpful than repeatedly querying the database
	if err := urlDataModel.LoadCache(); err != nil {
		return nil, err
	}

	return urlDataModel, nil
}

func (m *URLDataModel) LoadCache() error {
	stmt := `SELECT short_code, long_url, title, created_on, expires_on FROM urls`
	rows, err := m.db.Query(stmt)
	if err != nil {
		return err
	}

	for rows.Next() {
		var urlData models.URLData
		var expiryDate sql.NullTime
		err := rows.Scan(&urlData.ShortCode, &urlData.URL, &urlData.Title, &urlData.CreatedOn, &expiryDate)
		if err != nil {
			return err
		}
		if expiryDate.Valid {
			urlData.ExpiresOn = &expiryDate.Time
		}
		m.cache[urlData.ShortCode] = urlData
	}

	return rows.Err()
}

// returns all urls shortened from the database
func (m *URLDataModel) All() ([]models.URLData, error) {
	stmt := `
	SELECT short_code, long_url, title, created_on, expires_on
	FROM urls
	WHERE expires_on IS NULL 
	OR expires_on > datetime('now')
	ORDER BY created_on DESC;`

	rows, err := m.db.Query(stmt)
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
func (m *URLDataModel) CreateShortURL(url, title string, expiry time.Duration) (string, error) {
	shortCode := generateRandomString(m.shortURLLen)
	for {
		if _, exists := m.cache[shortCode]; !exists {
			break
		}
		shortCode = generateRandomString(m.shortURLLen)
	}

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

	_, err := m.db.Exec(stmt, urlData.ShortCode, urlData.URL, urlData.Title, urlData.CreatedOn, urlData.ExpiresOn)
	if err != nil {
		return "", err
	}

	// also update the cache
	m.cache[shortCode] = urlData

	return shortCode, nil
}

// fetches the data about the url from the cache
func (m *URLDataModel) GetRedirectData(shortCode string) (models.URLData, error) {
	urlData, exists := m.cache[shortCode]

	if !exists {
		return models.URLData{}, ErrNotExists
	}

	if urlData.ExpiresOn != nil && time.Now().After(*urlData.ExpiresOn) {
		return models.URLData{}, ErrURLExpired
	}

	return urlData, nil
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}
