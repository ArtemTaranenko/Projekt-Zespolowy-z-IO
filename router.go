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
	h := handlers{database: database}
	api := r.Group("/api")
	api.GET("/health", h.health)
	api.GET("/quiz/random", h.randomQuiz)
	// Kontrakty tras dla następnych etapów. Nie zgłaszają pozornego sukcesu.
	api.GET("/sentences/random", notImplemented)
	api.GET("/flashcards/random", notImplemented)
	api.POST("/flashcards", notImplemented)
	return r
}
