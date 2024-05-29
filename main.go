package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/restgoq")
	if err != nil {
		log.Fatal(err)
	}
	autoMigrate()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/create", insertData) // Correct route handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":8080", r)
}

func autoMigrate() {
	query := `
	CREATE TABLE IF NOT EXISTS usersss (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		age INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL DEFAULT NULL
	)`

	_, err = db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}

func insertData(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO usersss(name, age) VALUES(?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the prepared statement with the user data
	_, err = stmt.Exec(user.Name, user.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User successfully created"))
}
