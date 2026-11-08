package db

import (
	"context"
	"database/sql"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Word struct {
	ID      int      `json:"id"`
	English string   `json:"english"`
	Polish  string   `json:"polish"`
	Tags    []string `json:"tags"`
}

// initVocabulary tworzy brakujące tabele. Nie zmienia istniejących tabel i danych.
// Kolumny odpowiadają zapytaniom z dostarczonego main.go.
func initVocabulary(database *sql.DB) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	statements := []string{
		`CREATE TABLE IF NOT EXISTS tags (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY AUTOINCREMENT, english TEXT NOT NULL, polish TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS word_tags (word_id INTEGER NOT NULL, tag_id INTEGER NOT NULL, PRIMARY KEY (word_id, tag_id))`,
	}
	for _, statement := range statements {
		if _, err := tx.Exec(statement); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func ListTags(ctx context.Context, database *sql.DB) ([]Tag, error) {
	rows, err := database.QueryContext(ctx, `SELECT id, name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := []Tag{}
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// RandomWord odczytuje słowo i jego tagi w jednej transakcji odczytu.
func RandomWord(ctx context.Context, database *sql.DB, tagID int) (Word, error) {
	word := Word{Tags: []string{}}
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return word, err
	}
	defer tx.Rollback()
	err = tx.QueryRowContext(ctx, `SELECT w.id, w.english, w.polish
  FROM words w JOIN word_tags wt ON wt.word_id = w.id
  WHERE wt.tag_id = ? ORDER BY RANDOM() LIMIT 1`, tagID).Scan(&word.ID, &word.English, &word.Polish)
	if err != nil {
		return word, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT t.name FROM tags t
  JOIN word_tags wt ON wt.tag_id = t.id WHERE wt.word_id = ? ORDER BY t.name`, word.ID)
	if err != nil {
		return word, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return word, err
		}
		word.Tags = append(word.Tags, name)
	}
	if err := rows.Err(); err != nil {
		return word, err
	}
	if err := rows.Close(); err != nil {
		return word, err
	}
	return word, tx.Commit()
}
