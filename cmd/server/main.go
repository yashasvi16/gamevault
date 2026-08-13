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
	"context"
	"os/signal"
	"syscall"
	"time"
	"log/slog"
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
	
	slog.Info("Database connected succesfully")

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	
	playerRepo := repository.NewPlayerRepository(db)

	playerHandler := handler.NewPlayerHandler(playerRepo)
	authHandler := handler.NewAuthHandler(playerRepo)

	matchRepo := repository.NewMatchRepository(db)
	matchHandler := handler.NewMatchHandler(matchRepo)

	//Public routes
	r.Post("/players", playerHandler.RegisterPlayer)
	r.Get("/leaderboard", playerHandler.GetLeaderboard)
	r.Post("/login", authHandler.Login)
	r.Get("/health", handler.HealthCheck(db))

	//Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/profile", playerHandler.GetProfile)
		r.Post("/match", matchHandler.RecordMatch)
	})
	
	//http.ListenAndServe(":8080", r)

	srv := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	go func() {
		slog.Info("Server Starting", "port", 8080)
		if err := srv.ListenAndServe()
		err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx)
	err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	fmt.Println("Server exited gracefully")

}