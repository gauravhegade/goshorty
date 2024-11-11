package helpers

import (
	"context"
	"math/rand"
	"regexp"

	"github.com/jackc/pgx/v5"
)

// helper functions for the endpoints
func GenerateShortUrlKey(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"

	b := make([]byte, n)
	for i := range n {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// global map for storing shorturls
var ShortUrls = make(map[string]map[string]string)

// builds a map from individual components of each short url
// and assigns it to the global map variable
func BuildMap(key, longUrl, customAlias, creationDate, expiryDate string) {
	ShortUrls[key] = map[string]string{
		"url":          longUrl,
		"customAlias":  customAlias,
		"creationDate": creationDate,
		"expiryDate":   expiryDate,
	}
}

func IsValidDomainName(domainName string) bool {
	validNameRegex, _ := regexp.Compile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
	return validNameRegex.MatchString(domainName)
}

func NewPostgres(databaseURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
