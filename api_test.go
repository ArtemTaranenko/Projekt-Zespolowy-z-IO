package main

import (
	"database/sql"
	"encoding/json"
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
		{"GET", "/api/sentences/random", 404, "SENTENCE_NOT_FOUND"},
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
