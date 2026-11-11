package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// newRouter tworzy routing niezależnie od uruchomienia serwera HTTP.
func newRouter(database *sql.DB) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.HandleMethodNotAllowed = true
	r.NoRoute(func(c *gin.Context) { apiError(c, http.StatusNotFound, "NOT_FOUND", "Nie znaleziono endpointu.") })
	r.NoMethod(func(c *gin.Context) {
		apiError(c, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Niedozwolona metoda HTTP.")
	})
	r.GET("/", func(c *gin.Context) { c.File("./frontend/index.html") })
	r.GET("/style.css", func(c *gin.Context) { c.File("./frontend/style.css") })
	r.GET("/app.js", func(c *gin.Context) { c.File("./frontend/app.js") })
	r.GET("/quiz", func(c *gin.Context) { c.File("./frontend/quiz.html") })
	r.GET("/quiz.js", func(c *gin.Context) { c.File("./frontend/quiz.js") })
	r.GET("/quiz.css", func(c *gin.Context) { c.File("./frontend/quiz.css") })
	r.GET("/flashcards", func(c *gin.Context) { c.File("./frontend/flashcards.html") })
	r.GET("/sentences", func(c *gin.Context) { c.File("./frontend/sentences.html") })
	r.GET("/sentences.js", func(c *gin.Context) { c.File("./frontend/sentences.js") })
	r.GET("/sentences.css", func(c *gin.Context) { c.File("./frontend/sentences.css") })
	h := handlers{database: database}
	api := r.Group("/api")
	api.GET("/tags", h.tags)
	api.GET("/words/random", h.randomWord)
	api.GET("/health", h.health)
	api.GET("/quiz/random", h.randomQuiz)
	api.POST("/quiz/check", h.checkQuiz)
	// Kontrakty tras dla następnych etapów. Nie zgłaszają pozornego sukcesu.
	api.GET("/sentences/random", h.randomSentence)
	api.POST("/sentences/check", h.checkSentence)
	api.GET("/flashcards/random", notImplemented)
	api.POST("/flashcards", notImplemented)
	return r
}
