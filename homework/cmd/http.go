package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"test/internal/repository"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	db, err := sql.Open("mysql", user+":"+password+"@/users")
	if err != nil {
		log.Fatal("Error connecting DataBase")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Bad connection")
	}
	users := repository.CreateRepo(db)

	r := chi.NewRouter()

	r.Post("/create", users.Create)
	r.Post("/make_friends", users.MakeFriends)
	r.Put("/{id}", users.UpdateAge)
	r.Delete("/{id}", users.DeleteUser)
	r.Get("/{id}", users.GetFriends)


	log.Fatal(http.ListenAndServe(port, r))
}
