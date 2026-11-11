package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"english-learning-app/db"

	"github.com/gin-gonic/gin"
)

func (h handlers) randomQuiz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	quiz, err := db.RandomQuiz(ctx, h.database)
	if errors.Is(err, sql.ErrNoRows) {
		apiError(
			c,
			http.StatusNotFound,
			"QUIZ_NOT_FOUND",
			"Brak quizów w bazie danych.",
		)
		return
	}

	if err != nil {
		log.Printf("Odczyt quizu: %v", err)
		apiError(
			c,
			http.StatusInternalServerError,
			"DATABASE_ERROR",
			"Nie udało się pobrać quizu.",
		)
		return
	}

	answers := []string{
		quiz.CorrectAnswer,
		quiz.Wrong1,
		quiz.Wrong2,
		quiz.Wrong3,
	}

	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})

	c.JSON(http.StatusOK, gin.H{
		"id":       quiz.ID,
		"question": quiz.Question,
		"answers":  answers,
	})
}

func (h handlers) checkQuiz(c *gin.Context) {
	var input struct {
		QuizID int    `json:"quiz_id"`
		Answer string `json:"answer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		apiError(c, http.StatusBadRequest, "INVALID_JSON", "Nieprawidłowy JSON odpowiedzi.")
		return
	}

	input.Answer = strings.TrimSpace(input.Answer)
	if input.QuizID <= 0 || input.Answer == "" {
		apiError(c, http.StatusBadRequest, "INVALID_ANSWER", "Podaj identyfikator quizu i odpowiedź.")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	correctAnswer, err := db.QuizAnswer(ctx, h.database, input.QuizID)
	if errors.Is(err, sql.ErrNoRows) {
		apiError(c, http.StatusNotFound, "QUIZ_NOT_FOUND", "Nie znaleziono wskazanego quizu.")
		return
	}
	if err != nil {
		log.Printf("Sprawdzanie odpowiedzi quizu: %v", err)
		apiError(c, http.StatusInternalServerError, "DATABASE_ERROR", "Nie udało się sprawdzić odpowiedzi.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":        strings.EqualFold(input.Answer, correctAnswer),
		"correct_answer": correctAnswer,
	})
}
