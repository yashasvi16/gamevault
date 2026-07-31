package handler

import (
	"encoding/json"
	"net/http"
	"github.com/yashasvi16/gamevault/internal/repository"
	"github.com/yashasvi16/gamevault/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type PlayerHandler struct {
	playerRepo *repository.PlayerRepository
}

func NewPlayerHandler(repo *repository.PlayerRepository) *PlayerHandler {
	return &PlayerHandler {
		playerRepo: repo,
	}
}

func (h *PlayerHandler) RegisterPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newPlayer RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&newPlayer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error decoding request body",
		})
		return
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(newPlayer.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error generating password hash",
		})
		return
	}

	registeredPlayer := model.Player {
		Username: newPlayer.Username,
		Email: newPlayer.Email,
		PasswordHash: string(hashBytes),
	}

	err = h.playerRepo.CreatePlayer(&registeredPlayer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error registering new player",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any {
		"message": "Player registered successfully",
		"player": registeredPlayer,
	})
}