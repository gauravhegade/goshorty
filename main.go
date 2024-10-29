package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// new config parser
var k = koanf.New(".")

// global map for storing shorturls
var shorturls = make(map[string]map[string]string)

type shortUrl struct {
	LongUrl     string `json:"url"`
	CustomAlias string `json:"custom_alias"` // if any custom alias is given by the user
	ExpiryDate  string `json:"expiry_date"`  // if any expiry date is given by the user
}

func main() {
	log.Println("GO SHORTY IT'S YOUR BIRTHDAY")

	// get env vars
	if err := k.Load(file.Provider(".env.json"), json.Parser()); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read environment variables")
	}
	portNumber := k.String("app.port")

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies([]string{"*"})

	// api endpoints
	router.GET("/urls", getAllShortUrls)
	router.GET("/:key", getLongUrlByKey)
	router.POST("/urls", createShortUrl)

	router.Run(":" + portNumber)
}

func getAllShortUrls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, shorturls)
}

func createShortUrl(c *gin.Context) {
	var su shortUrl

	if err := c.BindJSON(&su); err != nil {
		c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "not a url"})
	}

	shortUrlKey := generateShortUrlKey(6)
	customFormat := time.DateOnly + " " + time.TimeOnly
	shortUrlCreationDate := time.Now().Format(customFormat)

	if validDomain := isValidDomainName(su.LongUrl); validDomain {
		buildMap(shortUrlKey, su.LongUrl, su.CustomAlias, shortUrlCreationDate, su.ExpiryDate)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "created short url: " + shortUrlKey})
		return
	}

	c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "domain name is not valid"})
}

func getLongUrlByKey(c *gin.Context) {
	key := c.Param("key")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "getting url by " + key})
}

// helper functions for the endpoints
func generateShortUrlKey(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"

	b := make([]byte, n)
	for i := range n {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// builds a map from individual components of each short url
// and assigns it to the global map variable
func buildMap(key, longUrl, customAlias, creationDate, expiryDate string) {
	shorturls[key] = map[string]string{
		"url":          longUrl,
		"customAlias":  customAlias,
		"creationDate": creationDate,
		"expiryDate":   expiryDate,
	}
}

func isValidDomainName(domainName string) bool {
	validNameRegex, _ := regexp.Compile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
	return validNameRegex.MatchString(domainName)
}
