package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"strings"
	"time"

	"english-learning-app/db"
	"github.com/gin-gonic/gin"
)

// sentenceWords usuwa końcową interpunkcję; zachowuje każde wystąpienie słowa.
func sentenceWords(english string) []string {
	return strings.Fields(strings.TrimRight(strings.TrimSpace(english), ".!?"))
}

func shuffledWords(english string) []string {
	original := sentenceWords(english)
	words := append([]string{}, original...)
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	// Jeśli losowanie pozostawiło poprawną kolejność, zamień dwie różne pozycje.
	if strings.Join(words, " ") == strings.Join(original, " ") {
		for i := 1; i < len(words); i++ {
			if !strings.EqualFold(words[0], words[i]) {
				words[0], words[i] = words[i], words[0]
				break
			}
		}
	}
	return words
}

func (h handlers) randomSentence(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	sentence, err := db.RandomSentence(ctx, h.database)
	if errors.Is(err, sql.ErrNoRows) {
		apiError(c, http.StatusNotFound, "SENTENCE_NOT_FOUND", "Brak zdań do ćwiczenia.")
		return
	}
	if err != nil {
		log.Printf("Losowanie zdania: %v", err)
		apiError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Nie udało się pobrać zdania.")
		return
	}
	words := shuffledWords(sentence.English)
	if len(words) < 2 || len(words) > 64 {
		apiError(c, http.StatusInternalServerError, "INVALID_EXERCISE", "Zdanie w bazie nie nadaje się do ćwiczenia.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": sentence.ID, "polish": sentence.Polish, "words": words})
}

func (h handlers) checkSentence(c *gin.Context) {
	mediaType, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
	if err != nil || mediaType != "application/json" {
		apiError(c, http.StatusUnsupportedMediaType, "INVALID_CONTENT_TYPE", "Wyślij odpowiedź w formacie JSON.")
		return
	}
	var input struct {
		SentenceID int      `json:"sentence_id"`
		Words      []string `json:"words"`
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8192)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&input)
	if err == nil {
		var extra any
		if trailingErr := decoder.Decode(&extra); trailingErr != io.EOF {
			if trailingErr != nil {
				err = trailingErr
			} else {
				err = errors.New("multiple JSON values")
			}
		}
	}
	if err != nil {
		var limit *http.MaxBytesError
		if errors.As(err, &limit) {
			apiError(c, http.StatusRequestEntityTooLarge, "BODY_TOO_LARGE", "Odpowiedź jest za długa.")
		} else {
			apiError(c, http.StatusBadRequest, "INVALID_JSON", "Nieprawidłowy JSON odpowiedzi.")
		}
		return
	}
	if input.SentenceID <= 0 || len(input.Words) == 0 || len(input.Words) > 64 {
		apiError(c, http.StatusBadRequest, "INVALID_ANSWER", "Podaj identyfikator zdania i listę słów.")
		return
	}
	for i, word := range input.Words {
		input.Words[i] = strings.TrimSpace(word)
		if len(input.Words[i]) == 0 || len(input.Words[i]) > 128 || len(strings.Fields(input.Words[i])) != 1 {
			apiError(c, http.StatusBadRequest, "INVALID_ANSWER", "Każdy element odpowiedzi musi być jednym słowem.")
			return
		}
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	sentence, err := db.SentenceByID(ctx, h.database, input.SentenceID)
	if errors.Is(err, sql.ErrNoRows) {
		apiError(c, http.StatusNotFound, "SENTENCE_NOT_FOUND", "Nie znaleziono wskazanego zdania.")
		return
	}
	if err != nil {
		log.Printf("Sprawdzanie zdania: %v", err)
		apiError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Nie udało się sprawdzić odpowiedzi.")
		return
	}
	expected := sentenceWords(sentence.English)
	correct := len(input.Words) == len(expected)
	if correct {
		for i := range expected {
			if !strings.EqualFold(input.Words[i], expected[i]) {
				correct = false
				break
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"correct": correct, "correct_answer": sentence.English})
}
