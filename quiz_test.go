package main

import (
	"encoding/json"
	"english-learning-app/db"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestRandomQuizUsesDatabase(t *testing.T) {
	database := testDB(t)
	_, err := database.Exec(`INSERT INTO quizzes (question, correct_answer, wrong1, wrong2, wrong3)
		VALUES (?, ?, ?, ?, ?)`, "Test from SQLite", "yes", "no", "maybe", "never")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	newRouter(database).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/quiz/random", nil))
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	var quiz struct {
		ID       int      `json:"id"`
		Question string   `json:"question"`
		Answers  []string `json:"answers"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &quiz); err != nil {
		t.Fatal(err)
	}

	if quiz.Question != "Test from SQLite" ||
		quiz.ID == 0 ||
		len(quiz.Answers) != 4 {
		t.Fatalf("wrong quiz: %+v", quiz)
	}

	var raw map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if _, exists := raw["correct_answer"]; exists {
		t.Fatal("random quiz exposes correct answer")
	}
}

func TestQuizSeedAndPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "test.db")
	database, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	for i := 0; i < 2; i++ {
		if err := db.SeedQuizzes(database); err != nil {
			t.Fatal(err)
		}
	}
	database.Close()
	reopened, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	var count int
	if err := reopened.QueryRow("SELECT COUNT(*) FROM quizzes").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 30 {
		t.Fatalf("expected 30 persistent rows, got %d", count)
	}
	if _, err := reopened.Exec("UPDATE quizzes SET question = 'Custom data'"); err != nil {
		t.Fatal(err)
	}
	if err := db.SeedQuizzes(reopened); err != nil {
		t.Fatal(err)
	}
	var question string
	if err := reopened.QueryRow("SELECT question FROM quizzes LIMIT 1").Scan(&question); err != nil {
		t.Fatal(err)
	}
	if question != "Custom data" {
		t.Fatal("seed overwrote existing data")
	}
}
