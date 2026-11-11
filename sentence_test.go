package main

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"english-learning-app/db"
)

func sentenceFixture(t *testing.T, database *sql.DB) {
	t.Helper()
	_, err := database.Exec(`INSERT INTO sentences (id, polish, english) VALUES (1, ?, ?)`,
		"Ten kot śpi na kanapie.", "The cat is sleeping on the sofa.")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSentenceRandomHidesAnswerAndKeepsDuplicateWords(t *testing.T) {
	database := testDB(t)
	sentenceFixture(t, database)
	router := newRouter(database)
	expected := []string{"The", "cat", "is", "sleeping", "on", "the", "sofa"}
	sorted := append([]string{}, expected...)
	slices.Sort(sorted)
	for i := 0; i < 20; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", "/api/sentences/random", nil))
		if w.Code != 200 {
			t.Fatalf("status %d: %s", w.Code, w.Body.String())
		}
		var body struct {
			ID     int      `json:"id"`
			Polish string   `json:"polish"`
			Words  []string `json:"words"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		var raw map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
			t.Fatal(err)
		}
		if len(raw) != 3 || raw["correct_answer"] != nil || raw["english"] != nil {
			t.Fatalf("answer exposed: %s", w.Body.String())
		}
		if body.ID != 1 || body.Polish != "Ten kot śpi na kanapie." {
			t.Fatalf("wrong exercise: %+v", body)
		}
		if slices.Equal(body.Words, expected) {
			t.Fatal("words not shuffled")
		}
		slices.Sort(body.Words)
		if !slices.Equal(body.Words, sorted) {
			t.Fatalf("missing or duplicated tokens: %v", body.Words)
		}
	}
	// Dwa identyczne słowa muszą zostać zachowane jako osobne kafelki.
	words := shuffledWords("I think that that is true.")
	count := 0
	for _, word := range words {
		if word == "that" {
			count++
		}
	}
	if count != 2 || len(words) != 6 {
		t.Fatalf("duplicates: %v", words)
	}
}

func TestSentenceCheck(t *testing.T) {
	database := testDB(t)
	sentenceFixture(t, database)
	router := newRouter(database)
	cases := []struct {
		name, body, contentType string
		status                  int
		correct                 bool
	}{
		{"correct", `{"sentence_id":1,"words":["The","cat","is","sleeping","on","the","sofa"]}`, "application/json", 200, true},
		{"case and whitespace", `{"sentence_id":1,"words":[" the ","CAT","is","sleeping","on","The","sofa"]}`, "application/json; charset=utf-8", 200, true},
		{"wrong order", `{"sentence_id":1,"words":["cat","The","is","sleeping","on","the","sofa"]}`, "application/json", 200, false},
		{"missing word", `{"sentence_id":1,"words":["The","cat"]}`, "application/json", 200, false},
		{"extra word", `{"sentence_id":1,"words":["The","cat","is","sleeping","on","the","sofa","today"]}`, "application/json", 200, false},
		{"unknown id", `{"sentence_id":999,"words":["hello"]}`, "application/json", 404, false},
		{"negative id", `{"sentence_id":-1,"words":["hello"]}`, "application/json", 400, false},
		{"empty list", `{"sentence_id":1,"words":[]}`, "application/json", 400, false},
		{"blank token", `{"sentence_id":1,"words":[" "]}`, "application/json", 400, false},
		{"two words in token", `{"sentence_id":1,"words":["The cat"]}`, "application/json", 400, false},
		{"string id", `{"sentence_id":"1 OR 1=1","words":["cat"]}`, "application/json", 400, false},
		{"invalid JSON", `{`, "application/json", 400, false},
		{"trailing JSON", `{"sentence_id":1,"words":["cat"]}{}`, "application/json", 400, false},
		{"client answer override", `{"sentence_id":1,"words":["cat"],"correct_answer":"cat"}`, "application/json", 400, false},
		{"too large", `{"sentence_id":1,"words":["` + strings.Repeat("x", 9000) + `"]}`, "application/json", 413, false},
		{"wrong content type", `{}`, "text/plain", 415, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/sentences/check", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			router.ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Fatalf("status %d, want %d: %s", w.Code, tc.status, w.Body.String())
			}
			if tc.status == 200 {
				var body struct {
					Correct       bool   `json:"correct"`
					CorrectAnswer string `json:"correct_answer"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.Correct != tc.correct || body.CorrectAnswer != "The cat is sleeping on the sofa." {
					t.Fatalf("result: %+v", body)
				}
			}
		})
	}
}

func TestSentenceSeedAndExistingData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.db")
	database, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if _, err := database.Exec(`INSERT INTO words (english,polish) VALUES ('cat','kot'); INSERT INTO tags(name) VALUES ('Animals'); INSERT INTO word_tags(word_id,tag_id) VALUES (1,1)`); err != nil {
		t.Fatal(err)
	}
	if err := db.SeedQuizzes(database); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := db.SeedSentences(database); err != nil {
			t.Fatal(err)
		}
	}
	database.Close()
	database, err = db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	for table, want := range map[string]int{"sentences": 10, "quizzes": 30, "words": 1, "tags": 1, "word_tags": 1} {
		var count int
		if err := database.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != want {
			t.Fatalf("%s count %d, want %d", table, count, want)
		}
	}
	if _, err := database.Exec(`UPDATE sentences SET english='Custom sentence.' WHERE id=1`); err != nil {
		t.Fatal(err)
	}
	if err := db.SeedSentences(database); err != nil {
		t.Fatal(err)
	}
	var english string
	if err := database.QueryRow(`SELECT english FROM sentences WHERE id=1`).Scan(&english); err != nil {
		t.Fatal(err)
	}
	if english != "Custom sentence." {
		t.Fatal("seed overwrote existing sentence")
	}
}

func TestSentenceDatabaseFailureAndStaticRoutes(t *testing.T) {
	database := testDB(t)
	router := newRouter(database)
	for _, path := range []string{"/sentences", "/sentences.js", "/sentences.css"} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 200 || w.Body.Len() == 0 {
			t.Fatalf("static %s: %d", path, w.Code)
		}
	}
	database.Close()
	for _, method := range []string{"GET", "POST"} {
		path := "/api/sentences/random"
		if method == "POST" {
			path = "/api/sentences/check"
		}
		req := httptest.NewRequest(method, path, strings.NewReader(`{"sentence_id":1,"words":["cat"]}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 500 {
			t.Fatalf("closed DB %s: %d", path, w.Code)
		}
	}
}
