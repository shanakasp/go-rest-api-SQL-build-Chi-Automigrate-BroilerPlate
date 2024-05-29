package main

import (
	"database/sql"
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
