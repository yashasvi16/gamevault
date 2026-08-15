package handler

import (
	"encoding/json"
	"net/http"
	"github.com/yashasvi16/gamevault/internal/repository"
	"github.com/yashasvi16/gamevault/internal/model"
	"github.com/yashasvi16/gamevault/internal/service"
	"github.com/yashasvi16/gamevault/internal/ws"
)

type MatchHandler struct {
	repo *repository.MatchRepository
	hub *ws.Hub
}

type RecordMatchRequest struct {
	OpponentID int `json:"opponent_id"`
	MyScore int `json:"my_score"`
	OpponentScore int `json:"opponent_score"`
}

func NewMatchHandler (repo *repository.MatchRepository, hub *ws.Hub) *MatchHandler {
	return &MatchHandler{
		repo: repo,
		hub: hub,
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

	winnerID := service.DetermineWinner(playerID, 
	req.OpponentID, req.MyScore, req.OpponentScore)

	var match model.Match
	match.Player1ID = playerID
	match.Player2ID = req.OpponentID
	match.WinnerID = winnerID
	match.Player1Score = req.MyScore
	match.Player2Score = req.OpponentScore

	err = h.repo.RecordMatchWithStats(&match)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error creating match",
		})
		return
	}

	// Broadcast to all WebSocket clients that the leaderboard has changed
	if h.hub != nil {
		notification := map[string]any{
			"type": "leaderboard_update",
			"message": "A match was just recorded",
		}
		notifJSON, err := json.Marshal(notification)
		if err == nil {
			h.hub.Broadcast(notifJSON)
		}
	}


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Match recorded successfully",
		"match": match,
	})

}