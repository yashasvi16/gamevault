package handler

import (
	"encoding/json"
	"net/http"
	"github.com/yashasvi16/gamevault/internal/repository"
	"github.com/yashasvi16/gamevault/internal/model"
)

type MatchHandler struct {
	repo *repository.MatchRepository
}

type RecordMatchRequest struct {
	OpponentID int `json:"opponent_id"`
	MyScore int `json:"my_score"`
	OpponentScore int `json:"opponent_score"`
}

func NewMatchHandler (repo *repository.MatchRepository) *MatchHandler {
	return &MatchHandler{
		repo: repo,
	}
}

func (h *MatchHandler) RecordMatch(w http.ResponseWriter, r *http.Request) {
	playerIDFloat := r.Context().Value("player_id").(float64)
	playerID := int(playerIDFloat)

	var req RecordMatchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error decoding request body",
		})
		return
	}

	var winnerID *int
	if req.MyScore > req.OpponentScore {
		winnerID = &playerID
	} else if req.MyScore < req.OpponentScore {
		winnerID = &req.OpponentID
	}

	var match model.Match
	match.Player1ID = playerID
	match.Player2ID = req.OpponentID
	match.WinnerID = winnerID
	match.Player1Score = req.MyScore
	match.Player2Score = req.OpponentScore

	err = h.repo.CreateMatch(&match)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error creating match",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Match recorded successfully",
		"match": match,
	})

}