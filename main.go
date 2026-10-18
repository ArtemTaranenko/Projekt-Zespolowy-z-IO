package main

import (
	"database/sql"
	"log"

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
