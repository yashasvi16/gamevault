package handler

import (
	"encoding/json"
	"net/http"
	"github.com/yashasvi16/gamevault/internal/repository"
	"github.com/yashasvi16/gamevault/internal/model"
	"github.com/yashasvi16/gamevault/internal/service"
	"github.com/yashasvi16/gamevault/internal/worker"
)

type MatchHandler struct {
	repo *repository.MatchRepository
	statsJobs chan worker.StatsJob
}

type RecordMatchRequest struct {
	OpponentID int `json:"opponent_id"`
	MyScore int `json:"my_score"`
	OpponentScore int `json:"opponent_score"`
}

func NewMatchHandler (repo *repository.MatchRepository, statsJobs chan worker.StatsJob) *MatchHandler {
	return &MatchHandler{
		repo: repo,
		statsJobs: statsJobs,
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

	err = h.repo.CreateMatch(&match)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error creating match",
		})
		return
	}

	h.statsJobs <- worker.StatsJob{
		PlayerID: match.Player1ID,
		Won: match.WinnerID != nil && 
			*match.WinnerID == match.Player1ID,
	}
	h.statsJobs <- worker.StatsJob{
		PlayerID: match.Player2ID,
		Won: match.WinnerID != nil &&
			*match.WinnerID == match.Player2ID,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Match recorded successfully",
		"match": match,
	})

}