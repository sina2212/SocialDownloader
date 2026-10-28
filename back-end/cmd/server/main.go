package main

import (
	"back-end/internal/api"
	"back-end/internal/downloader"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	dl := downloader.New("./downloads")

	handler := api.NewHandler(dl)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/extract", handler.Extract)
	r.POST("/download", handler.Download)

	log.Println("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}