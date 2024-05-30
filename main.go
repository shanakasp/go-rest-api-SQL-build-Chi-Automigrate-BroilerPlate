package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int    `json:"id"`
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

	r.Post("/create", insertData) 
	r.Get("/getallusers", getAll) 
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

	responseMessage := fmt.Sprintf("User successfully created: %s, Age: %d", user.Name, user.Age)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseMessage))
}
func getAll(w http.ResponseWriter, r *http.Request) {
	// Prepare the SQL statement
	rows, err := db.Query("SELECT id, name, age FROM usersss")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	// Iterate over the rows
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the users slice to JSON
	responseData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

