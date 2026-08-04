package handler

import (
	"encoding/json"
	"net/http"
	"github.com/yashasvi16/gamevault/internal/repository"
	"github.com/yashasvi16/gamevault/internal/model"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type PlayerHandler struct {
	playerRepo repository.PlayerRepo
}

func NewPlayerHandler(repo repository.PlayerRepo) *PlayerHandler {
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

	err = validate.Struct(newPlayer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Validation failed: " + err.Error(),
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

func (h *PlayerHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset := (page - 1) * limit

	players, err := h.playerRepo.GetLeaderboard(limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string {
			"message": "Error fetching leaderboard",
		})
		return
	}

	w.WriteHeader(http.StatusOK) 
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Leaderboard fetched successfully",
		"data": players,
	})
}

func (h *PlayerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	playerID := r.Context().Value("player_id")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Welcome to your dashboard",
		"player_id": playerID,
	})
}