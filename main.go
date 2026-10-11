package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/", "./frontend")

	api := r.Group("/api")
	{
		// Osoba 1: Moduł układania zdań
		api.GET("/sentences/random", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Endpoint dla zdań gotowy"})
		})

		// Osoba 2: Moduł fiszek
		api.GET("/flashcards/random", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Endpoint pobierania fiszek gotowy"})
		})
		api.POST("/flashcards", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Endpoint dodawania fiszek gotowy"})
		})

		// Osoba 3: Moduł quizu
		api.GET("/quiz/random", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Endpoint quizu gotowy"})
		})
	}

	r.Run(":8080")
}
