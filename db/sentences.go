package db

import (
	"context"
	"database/sql"
)

// Sentence przechowuje wzorcową odpowiedź wyłącznie po stronie serwera.
type Sentence struct {
	ID      int
	Polish  string
	English string
}

func initSentences(database *sql.DB) error {
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS sentences (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  polish TEXT NOT NULL,
  english TEXT NOT NULL
 )`)
	return err
}

func RandomSentence(ctx context.Context, database *sql.DB) (Sentence, error) {
	var sentence Sentence
	err := database.QueryRowContext(ctx, `SELECT id, polish, english FROM sentences
  ORDER BY RANDOM() LIMIT 1`).Scan(&sentence.ID, &sentence.Polish, &sentence.English)
	return sentence, err
}

func SentenceByID(ctx context.Context, database *sql.DB, id int) (Sentence, error) {
	var sentence Sentence
	err := database.QueryRowContext(ctx, `SELECT id, polish, english FROM sentences WHERE id = ?`, id).
		Scan(&sentence.ID, &sentence.Polish, &sentence.English)
	return sentence, err
}

// SeedSentences dodaje przykłady tylko do pustej tabeli, bez zmiany słówek i quizów.
func SeedSentences(database *sql.DB) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM sentences`).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	examples := []Sentence{
		{Polish: "Codziennie piję wodę.", English: "I drink water every day."},
		{Polish: "Ona czyta książkę.", English: "She is reading a book."},
		{Polish: "Mój brat lubi muzykę.", English: "My brother likes music."},
		{Polish: "Nie lubię zimnej kawy.", English: "I do not like cold coffee."},
		{Polish: "Ten kot śpi na kanapie.", English: "The cat is sleeping on the sofa."},
		{Polish: "Oni uczą się angielskiego.", English: "They are learning English."},
		{Polish: "Gdzie jest twój plecak?", English: "Where is your backpack?"},
		{Polish: "Wczoraj kupiliśmy nowy stół.", English: "We bought a new table yesterday."},
		{Polish: "Mam jabłko i banana.", English: "I have an apple and a banana."},
		{Polish: "Czy możesz mi pomóc?", English: "Can you help me?"},
	}
	for _, sentence := range examples {
		if _, err := tx.Exec(`INSERT INTO sentences (polish, english) VALUES (?, ?)`, sentence.Polish, sentence.English); err != nil {
			return err
		}
	}
	return tx.Commit()
}
