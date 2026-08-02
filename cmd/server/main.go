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
	"github.com/yashasvi16/gamevault/internal/middleware"
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

	playerRepo := repository.NewPlayerRepository(db)
	playerHandler := handler.NewPlayerHandler(playerRepo)
	authHandler := handler.NewAuthHandler(playerRepo)

	matchRepo := repository.NewMatchRepository(db)
	matchHandler := handler.NewMatchHandler(matchRepo)

	//Public routes
	r.Post("/players", playerHandler.RegisterPlayer)
	r.Get("/leaderboard", playerHandler.GetLeaderboard)
	r.Post("/login", authHandler.Login)

	//Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/profile", playerHandler.GetProfile)
		r.Post("/match", matchHandler.RecordMatch)
	})
	
	http.ListenAndServe(":8080", r)

}