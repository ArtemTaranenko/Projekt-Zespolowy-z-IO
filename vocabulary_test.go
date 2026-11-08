package main

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"english-learning-app/db"
)

func TestVocabularyValidationAndEmptyData(t *testing.T) {
	database := testDB(t)
	router := newRouter(database)
	cases := []struct {
		path   string
		status int
	}{
		{"/api/tags", 200},
		{"/api/words/random", 400},
		{"/api/words/random?tag_id=abc", 400},
		{"/api/words/random?tag_id=0", 400},
		{"/api/words/random?tag_id=-1", 400},
		{"/api/words/random?tag_id=999999999999999999999", 400},
		{"/api/words/random?tag_id=1%20OR%201=1", 400},
		{"/api/words/random?tag_id=999", 404},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status %d: %s", w.Code, w.Body.String())
			}
			if tc.path == "/api/tags" {
				if w.Body.String() != "[]" {
					t.Fatalf("expected empty JSON array, got %s", w.Body.String())
				}
			} else {
				var body struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.Error == "" {
					t.Fatal("missing string error")
				}
			}
		})
	}
}

func TestVocabularySelectionAndPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.db")
	database, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	statements := []string{
		`INSERT INTO tags (id, name) VALUES (1, 'Animals'), (2, 'Basics'), (3, 'Food')`,
		`INSERT INTO words (id, english, polish) VALUES (1, 'cat', 'kot'), (2, 'apple', 'jabłko')`,
		`INSERT INTO word_tags (word_id, tag_id) VALUES (1, 1), (1, 2), (2, 3)`,
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	database.Close()
	database, err = db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	router := newRouter(database)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/api/tags", nil))
	if w.Code != 200 {
		t.Fatalf("tags status %d", w.Code)
	}
	var tags []db.Tag
	if err := json.Unmarshal(w.Body.Bytes(), &tags); err != nil {
		t.Fatal(err)
	}
	if len(tags) != 3 || tags[0].Name != "Animals" || tags[2].Name != "Food" {
		t.Fatalf("tags: %+v", tags)
	}
	for _, tc := range []struct {
		id, english string
		tagCount    int
	}{{"1", "cat", 2}, {"3", "apple", 1}} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", "/api/words/random?tag_id="+tc.id, nil))
		if w.Code != 200 {
			t.Fatalf("word status %d: %s", w.Code, w.Body.String())
		}
		var word db.Word
		if err := json.Unmarshal(w.Body.Bytes(), &word); err != nil {
			t.Fatal(err)
		}
		if word.English != tc.english || len(word.Tags) != tc.tagCount {
			t.Fatalf("word: %+v", word)
		}
		if word.English == "cat" && (word.Polish != "kot" || word.Tags[0] != "Animals" || word.Tags[1] != "Basics") {
			t.Fatalf("tags: %+v", word)
		}
	}
}

func TestVocabularyDatabaseFailure(t *testing.T) {
	database := testDB(t)
	database.Close()
	router := newRouter(database)
	for _, path := range []string{"/api/tags", "/api/words/random?tag_id=1"} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 500 {
			t.Fatalf("%s: %d", path, w.Code)
		}
		var body struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Error == "" {
			t.Fatal("missing error")
		}
	}
}

func TestNewPageRoutes(t *testing.T) {
	database := testDB(t)
	// Tymczasowe pliki sprawdzają routing, nie zastępują frontendu kolegi.
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	temp := t.TempDir()
	if err := os.Mkdir(filepath.Join(temp, "frontend"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"quiz", "flashcards"} {
		if err := os.WriteFile(filepath.Join(temp, "frontend", name+".html"), []byte(name), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chdir(temp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(original)
	router := newRouter(database)
	for _, name := range []string{"quiz", "flashcards"} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", "/"+name, nil))
		if w.Code != 200 || w.Body.String() != name {
			t.Fatalf("%s: %d %s", name, w.Code, w.Body.String())
		}
	}
}
