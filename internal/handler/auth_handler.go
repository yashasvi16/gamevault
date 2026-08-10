package handler

import (
	"encoding/json"
	"net/http"
	"github.com/yashasvi16/gamevault/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"github.com/yashasvi16/gamevault/internal/model"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"errors"
)

type AuthHandler struct {
	repo repository.PlayerRepo
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(repo repository.PlayerRepo) *AuthHandler {
	return &AuthHandler{
		repo: repo,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error decoding request body",
		})	
		return
	}

	var player *model.Player
	player, err = h.repo.GetPlayerByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid email or password",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.PasswordHash), []byte(req.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"player_id": player.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"message": "Login successful",
	})
}