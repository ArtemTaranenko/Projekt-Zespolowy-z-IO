package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"english-learning-app/db"
	"github.com/gin-gonic/gin"
)

// Nowe endpointy zachowują error jako tekst, zgodnie z wersją kolegi.
// Nie znamy kodu nowego frontendu, więc nie zmieniamy jego kontraktu.
func (h handlers) tags(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	tags, err := db.ListTags(ctx, h.database)
	if err != nil {
		log.Printf("Odczyt tagów: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie udało się pobrać tagów"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func (h handlers) randomWord(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Query("tag_id"))
	if err != nil || tagID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wybierz prawidłowy tag"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	word, err := db.RandomWord(ctx, h.database, tagID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nie znaleziono słówek dla wybranego tagu"})
		return
	}
	if err != nil {
		log.Printf("Losowanie słówka i odczyt tagów: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie udało się wylosować słówka"})
		return
	}
	c.JSON(http.StatusOK, word)
}
