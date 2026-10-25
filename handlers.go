package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"english-learning-app/db"
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

func (h handlers) randomQuiz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	quiz, err := db.RandomQuiz(ctx, h.database)
	if errors.Is(err, sql.ErrNoRows) {
		apiError(c, http.StatusNotFound, "QUIZ_NOT_FOUND", "Brak quizów w bazie danych.")
		return
	}
	if err != nil {
		log.Printf("Odczyt quizu: %v", err)
		apiError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Nie udało się pobrać quizu.")
		return
	}
	c.JSON(http.StatusOK, quiz)
}

func notImplemented(c *gin.Context) {
	apiError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Moduł zostanie zrealizowany w kolejnym etapie projektu.")
}
