package db

import (
	"context"
	"database/sql"
)

type Quiz struct {
	ID            int    `json:"id"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Wrong1        string `json:"wrong1"`
	Wrong2        string `json:"wrong2"`
	Wrong3        string `json:"wrong3"`
}

func initQuizzes(database *sql.DB) error {
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS quizzes (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        question TEXT,
        correct_answer TEXT,
        wrong1 TEXT,
        wrong2 TEXT,
        wrong3 TEXT
    )`)
	return err
}

func SeedQuizzes(database *sql.DB) error {
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
		{Question: "They ____ dinner now.", CorrectAnswer: "are having", Wrong1: "have", Wrong2: "had", Wrong3: "having"},
		{Question: "I ____ this book last week.", CorrectAnswer: "read", Wrong1: "reads", Wrong2: "am reading", Wrong3: "reading"},
		{Question: "He ____ football on Sundays.", CorrectAnswer: "plays", Wrong1: "play", Wrong2: "played", Wrong3: "playing"},
		{Question: "We ____ to London tomorrow.", CorrectAnswer: "are going", Wrong1: "go", Wrong2: "went", Wrong3: "going"},
		{Question: "She has ____ her homework.", CorrectAnswer: "finished", Wrong1: "finish", Wrong2: "finishing", Wrong3: "finishes"},
		{Question: "There ____ two books on the table.", CorrectAnswer: "are", Wrong1: "is", Wrong2: "be", Wrong3: "was"},
		{Question: "Tom ____ coffee every morning.", CorrectAnswer: "drinks", Wrong1: "drink", Wrong2: "drank", Wrong3: "drinking"},
		{Question: "My sister ____ TV at the moment.", CorrectAnswer: "is watching", Wrong1: "watches", Wrong2: "watched", Wrong3: "watch"},
		{Question: "They ____ in this house since 2015.", CorrectAnswer: "have lived", Wrong1: "live", Wrong2: "lived", Wrong3: "are living"},
		{Question: "He ____ his keys yesterday.", CorrectAnswer: "lost", Wrong1: "loses", Wrong2: "has lost", Wrong3: "losing"},
		{Question: "I ____ breakfast before school every day.", CorrectAnswer: "eat", Wrong1: "eats", Wrong2: "ate", Wrong3: "eating"},
		{Question: "Look! It ____ outside.", CorrectAnswer: "is raining", Wrong1: "rains", Wrong2: "rained", Wrong3: "rain"},
		{Question: "We ____ our grandparents last weekend.", CorrectAnswer: "visited", Wrong1: "visit", Wrong2: "visiting", Wrong3: "visits"},
		{Question: "She ____ never ____ sushi.", CorrectAnswer: "has, eaten", Wrong1: "have, eat", Wrong2: "is, eating", Wrong3: "did, ate"},
		{Question: "John ____ his homework yet.", CorrectAnswer: "hasn't finished", Wrong1: "doesn't finish", Wrong2: "didn't finished", Wrong3: "isn't finish"},
		{Question: "My parents ____ to work by car.", CorrectAnswer: "go", Wrong1: "goes", Wrong2: "going", Wrong3: "went"},
		{Question: "The children ____ in the park yesterday.", CorrectAnswer: "played", Wrong1: "play", Wrong2: "plays", Wrong3: "playing"},
		{Question: "I think it ____ tomorrow.", CorrectAnswer: "will rain", Wrong1: "rains", Wrong2: "rained", Wrong3: "raining"},
		{Question: "She ____ a letter right now.", CorrectAnswer: "is writing", Wrong1: "writes", Wrong2: "wrote", Wrong3: "write"},
		{Question: "We ____ this movie before.", CorrectAnswer: "have seen", Wrong1: "see", Wrong2: "saw", Wrong3: "are seeing"},
		{Question: "Peter ____ his room every Saturday.", CorrectAnswer: "cleans", Wrong1: "clean", Wrong2: "cleaned", Wrong3: "cleaning"},
		{Question: "They ____ football when it started to rain.", CorrectAnswer: "were playing", Wrong1: "played", Wrong2: "are playing", Wrong3: "play"},
		{Question: "I ____ my friend at six o'clock tomorrow.", CorrectAnswer: "will meet", Wrong1: "meet", Wrong2: "met", Wrong3: "meeting"},
		{Question: "There ____ a lot of people at the concert.", CorrectAnswer: "were", Wrong1: "was", Wrong2: "are", Wrong3: "be"},
		{Question: "You ____ wear a seat belt in a car.", CorrectAnswer: "must", Wrong1: "can", Wrong2: "might", Wrong3: "would"},
		{Question: "Would you like ____ tea?", CorrectAnswer: "some", Wrong1: "any", Wrong2: "many", Wrong3: "much"},
		{Question: "This exercise is ____ than the previous one.", CorrectAnswer: "easier", Wrong1: "easy", Wrong2: "more easy", Wrong3: "easiest"},
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

func RandomQuiz(ctx context.Context, database *sql.DB) (Quiz, error) {
	var quiz Quiz

	err := database.QueryRowContext(ctx, `
        SELECT id, question, correct_answer, wrong1, wrong2, wrong3
        FROM quizzes
        ORDER BY RANDOM()
        LIMIT 1
    `).Scan(
		&quiz.ID,
		&quiz.Question,
		&quiz.CorrectAnswer,
		&quiz.Wrong1,
		&quiz.Wrong2,
		&quiz.Wrong3,
	)

	return quiz, err
}

func QuizAnswer(ctx context.Context, database *sql.DB, id int) (string, error) {
	var answer string

	err := database.QueryRowContext(ctx,
		`SELECT correct_answer FROM quizzes WHERE id = ?`,
		id,
	).Scan(&answer)

	return answer, err
}
