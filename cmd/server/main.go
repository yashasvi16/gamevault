package main

import (
	"log"
	"os"
	"database/sql"
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/yashasvi16/gamevault/internal/repository"
	"github.com/yashasvi16/gamevault/internal/handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting Database: ", err)
	}
	
	fmt.Println("Database connected succesfully")

	r := chi.NewRouter()

	repo := repository.NewPlayerRepository(db)
	playerHandler := handler.NewPlayerHandler(repo)
	r.Post("/players", playerHandler.RegisterPlayer)

	authHandler := handler.NewAuthHandler(repo)
	r.Post("/login", authHandler.Login)
	
	http.ListenAndServe(":8080", r)

}