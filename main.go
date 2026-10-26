package main

import (
	"database/sql"
	"log"

    "net/http"
    "strconv"


	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Quiz struct {
	ID            int    `json:"id"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Wrong1        string `json:"wrong1"`
	Wrong2        string `json:"wrong2"`
	Wrong3        string `json:"wrong3"`
}

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


var db *sql.DB

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
	r.GET("/style.css", func(c *gin.Context) {
		c.File("./frontend/style.css")
	})
	r.GET("/app.js", func(c *gin.Context) {
		c.File("./frontend/app.js")
	})

	r.GET("/quiz", func(c *gin.Context) {
		c.File("./frontend/quiz.html")
	})

	r.GET("/flashcards", func(c *gin.Context) {
		c.File("./frontend/flashcards.html")
	})


	api := r.Group("/api")
	{
		// Osoba 1: Moduł układania zdań
		api.GET("/sentences/random", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Endpoint dla zdań gotowy"})
		})

		// Osoba 2: Moduł fiszek
		api.GET("/flashcards/random", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Endpoint pobierania fiszek gotowy"})
		})
		api.POST("/flashcards", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Endpoint dodawania fiszek gotowy"})
		})

		// KP: Moduł quizu
		api.GET("/quiz/random", func(c *gin.Context) {
			// Testowy obiekt
			exampleQuiz := Quiz{
				ID:            1,
				Question:      "Look! The dog ____ in the garden.",
				CorrectAnswer: "is playing",
				Wrong1:        "plays",
				Wrong2:        "played",
				Wrong3:        "play",
			}
			c.JSON(200, exampleQuiz)
		})
	}

r.GET("/api/tags", func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT id, name
			FROM tags
			ORDER BY name
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Nie udało się pobrać tagów",
			})
			return
		}
		defer rows.Close()

		tags := []Tag{}

		for rows.Next() {
			var tag Tag
			if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Nie udało się odczytać tagów",
				})
				return
			}
			tags = append(tags, tag)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Błąd podczas odczytu tagów",
			})
			return
		}

		c.JSON(http.StatusOK, tags)
	})

	// Losowanie słówka przypisanego do wybranego tagu
	r.GET("/api/words/random", func(c *gin.Context) {
		tagID, err := strconv.Atoi(c.Query("tag_id"))
		if err != nil || tagID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Wybierz prawidłowy tag",
			})
			return
		}

		var word Word

		err = db.QueryRow(`
			SELECT w.id, w.english, w.polish
			FROM words w
			JOIN word_tags wt ON wt.word_id = w.id
			WHERE wt.tag_id = ?
			ORDER BY RANDOM()
			LIMIT 1
		`, tagID).Scan(&word.ID, &word.English, &word.Polish)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Nie znaleziono słówek dla wybranego tagu",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Nie udało się wylosować słówka",
			})
			return
		}

		// Pobierz wszystkie tagi wylosowanego słówka
		rows, err := db.Query(`
			SELECT t.name
			FROM tags t
			JOIN word_tags wt ON wt.tag_id = t.id
			WHERE wt.word_id = ?
			ORDER BY t.name
		`, word.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Nie udało się pobrać tagów słówka",
			})
			return
		}
		defer rows.Close()

		word.Tags = []string{}

		for rows.Next() {
			var tagName string
			if err := rows.Scan(&tagName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Nie udało się odczytać tagów słówka",
				})
				return
			}
			word.Tags = append(word.Tags, tagName)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Błąd podczas odczytu tagów słówka",
			})
			return
		}

		c.JSON(http.StatusOK, word)
	})

	
	r.Run(":8080")
}

func initDB() {
	var err error

	db, err = sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		log.Fatal(err)
	}

	createQuizTableSQL := `CREATE TABLE IF NOT EXISTS quizzes (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"question" TEXT,
		"correct_answer" TEXT,
		"wrong1" TEXT,
		"wrong2" TEXT,
		"wrong3" TEXT
	);`

	_, err = db.Exec(createQuizTableSQL)
	if err != nil {
		log.Fatal("Nie udało się utworzyć tabeli quizzes: ", err)
	}
	log.Println("Baza danych i tabela dla quizu zostały poprawnie załadowane/utworzone!")
}
