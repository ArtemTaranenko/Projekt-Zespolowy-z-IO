package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"english-learning-app/db"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

func TestRoutes(t *testing.T) {
	database := testDB(t)
	router := newRouter(database)
	cases := []struct {
		method, path string
		status       int
		code         string
	}{
		{"GET", "/api/health", 200, ""},
		{"GET", "/api/quiz/random", 404, "QUIZ_NOT_FOUND"},
		{"GET", "/api/sentences/random", 501, "NOT_IMPLEMENTED"},
		{"GET", "/api/flashcards/random", 501, "NOT_IMPLEMENTED"},
		{"POST", "/api/flashcards", 501, "NOT_IMPLEMENTED"},
		{"POST", "/api/quiz/random", 405, "METHOD_NOT_ALLOWED"},
		{"GET", "/api/missing", 404, "NOT_FOUND"},
	}
	for _, tc := range cases {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status %d: %s", w.Code, w.Body.String())
			}
			var body struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Error.Code != tc.code {
				t.Fatalf("error code: %s", body.Error.Code)
			}
		})
	}
}

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
	var quiz db.Quiz
	if err := json.Unmarshal(w.Body.Bytes(), &quiz); err != nil {
		t.Fatal(err)
	}
	if quiz.Question != "Test from SQLite" || quiz.CorrectAnswer != "yes" || quiz.ID == 0 {
		t.Fatalf("wrong quiz: %+v", quiz)
	}
}

func TestSeedAndPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "test.db")
	database, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	for i := 0; i < 2; i++ {
		if err := db.Seed(database); err != nil {
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
	if count != 3 {
		t.Fatalf("expected 3 persistent rows, got %d", count)
	}
	if _, err := reopened.Exec("UPDATE quizzes SET question = 'Custom data'"); err != nil {
		t.Fatal(err)
	}
	if err := db.Seed(reopened); err != nil {
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

func TestDatabaseFailure(t *testing.T) {
	database := testDB(t)
	database.Close()
	router := newRouter(database)
	for _, tc := range []struct {
		path   string
		status int
	}{{"/api/health", 503}, {"/api/quiz/random", 500}} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != tc.status {
			t.Fatalf("%s: %d", tc.path, w.Code)
		}
	}
}
