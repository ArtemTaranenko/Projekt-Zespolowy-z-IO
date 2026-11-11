// Package db odpowiada za połączenie z SQLite i odczyt danych.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

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
		err = initQuizzes(database)
	}
	if err == nil {
		err = initVocabulary(database)
	}
	if err == nil {
		err = initSentences(database)
	}
	if err != nil {
		database.Close()
		return nil, fmt.Errorf("inicjalizacja SQLite: %w", err)
	}
	return database, nil
}
