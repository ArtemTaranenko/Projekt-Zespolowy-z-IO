package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type handlers struct{ database *sql.DB }

func apiError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

func (h handlers) health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := h.database.PingContext(ctx); err != nil {
		log.Printf("Sprawdzenie SQLite: %v", err)
		apiError(c, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "Baza danych jest niedostępna.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "database": "ok"})
}

func notImplemented(c *gin.Context) {
	apiError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Moduł zostanie zrealizowany w kolejnym etapie projektu.")
}
