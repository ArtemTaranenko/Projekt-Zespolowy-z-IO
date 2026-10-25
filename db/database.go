// Package db odpowiada za połączenie z SQLite i odczyt danych.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Quiz zachowuje format JSON z pierwszej wersji projektu.
type Quiz struct {
	ID            int    `json:"id"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Wrong1        string `json:"wrong1"`
	Wrong2        string `json:"wrong2"`
	Wrong3        string `json:"wrong3"`
}

// Open otwiera bazę i tworzy tabelę, nie usuwając istniejących danych.
func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(1)
	if err = database.Ping(); err == nil {
		_, err = database.Exec(`CREATE TABLE IF NOT EXISTS quizzes (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			question TEXT, correct_answer TEXT,
			wrong1 TEXT, wrong2 TEXT, wrong3 TEXT
		)`)
	}
	if err != nil {
		database.Close()
		return nil, fmt.Errorf("inicjalizacja SQLite: %w", err)
	}
	return database, nil
}

// RandomQuiz zwraca sql.ErrNoRows, gdy tabela jest pusta.
func RandomQuiz(ctx context.Context, database *sql.DB) (Quiz, error) {
	var q Quiz
	err := database.QueryRowContext(ctx, `SELECT id, question, correct_answer,
		wrong1, wrong2, wrong3 FROM quizzes ORDER BY RANDOM() LIMIT 1`).Scan(
		&q.ID, &q.Question, &q.CorrectAnswer, &q.Wrong1, &q.Wrong2, &q.Wrong3)
	return q, err
}

// Seed dodaje dane demonstracyjne wyłącznie do pustej tabeli.
// Transakcja zapewnia, że wszystkie wpisy zostaną dodane razem.
func Seed(database *sql.DB) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int
	if err = tx.QueryRow("SELECT COUNT(*) FROM quizzes").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	examples := []Quiz{
		{Question: "Look! The dog ____ in the garden.", CorrectAnswer: "is playing", Wrong1: "plays", Wrong2: "played", Wrong3: "play"},
		{Question: "She ____ to school every day.", CorrectAnswer: "goes", Wrong1: "go", Wrong2: "going", Wrong3: "gone"},
		{Question: "Yesterday we ____ a film.", CorrectAnswer: "watched", Wrong1: "watch", Wrong2: "watching", Wrong3: "watches"},
	}
	for _, q := range examples {
		_, err = tx.Exec(`INSERT INTO quizzes (question, correct_answer, wrong1, wrong2, wrong3)
			VALUES (?, ?, ?, ?, ?)`, q.Question, q.CorrectAnswer, q.Wrong1, q.Wrong2, q.Wrong3)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
