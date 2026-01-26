package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.GET("/query", func(c *gin.Context) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", os.Getenv("HYTALE_QUERY_URL"), nil)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		req.SetBasicAuth(os.Getenv("HYTALE_QUERY_USERNAME"), os.Getenv("HYTALE_QUERY_PASSWORD"))

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var result interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(resp.StatusCode, result)
	})

	r.Run(":8080")
}
