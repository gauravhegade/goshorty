package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gauravhegade/goshorty/helpers"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// new config parser
var k = koanf.New(".")

type shortUrl struct {
	LongUrl     string `json:"url" binding:"required"`
	CustomAlias string `json:"custom_alias"` // optional
	ExpiryDate  string `json:"expiry_date"`  // optional
}

func main() {
	log.Println("GO SHORTY IT'S YOUR BIRTHDAY")

	// get env vars
	if err := k.Load(file.Provider(".env.json"), json.Parser()); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read environment variables")
	}

	// connect to database
	databaseURL := k.String("database_url")
	conn, err := helpers.NewPostgres(databaseURL)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
		return
	}
	defer conn.Close(context.Background())

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies([]string{"*"})

	// api endpoints
	router.GET("/", homeHandler)
	router.GET("/urls", getAllShortUrls)
	router.GET("/:key", getLongUrlByKey)
	router.POST("/urls", createShortUrl)

	// get port number from env
	portNumber := k.String("port")
	router.Run(":" + portNumber)
}

func homeHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("URL Shortener Service"))
}

func getAllShortUrls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, helpers.ShortUrls)
}

func createShortUrl(c *gin.Context) {
	var su shortUrl

	if err := c.ShouldBind(&su); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrlKey := helpers.GenerateShortUrlKey(6)
	t := time.Now()
	customTimeFormat := "2006-01-02 15:04:05-07:00" // custom date format to adhere to postgres timestamptz format
	shortUrlCreationDate := t.Format(customTimeFormat)

	shortUrlExpiryDate, _ := time.Parse(customTimeFormat, su.ExpiryDate)
	parsedCreationDate, _ := time.Parse(customTimeFormat, shortUrlCreationDate)
	if shortUrlExpiryDate.Before(parsedCreationDate) {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Expiry date should be a future date"})
		return
	}

	if validDomain := helpers.IsValidDomainName(su.LongUrl); validDomain {
		helpers.BuildMap(shortUrlKey, su.LongUrl, su.CustomAlias, shortUrlCreationDate, su.ExpiryDate)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "created short url: " + shortUrlKey})
		return
	}

	c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "domain name is not valid"})
}

func getLongUrlByKey(c *gin.Context) {
	key := c.Param("key")

	if value, ok := helpers.ShortUrls[key]; ok {
		message := fmt.Sprintf(
			"Fetched URL: %s, with alias: %s, creation date: %s, expiry date: %s.",
			value["url"],
			value["customAlias"],
			value["creationDate"],
			value["expiryDate"],
		)
		c.IndentedJSON(http.StatusOK, gin.H{"message": message})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "could not find key"})
}
