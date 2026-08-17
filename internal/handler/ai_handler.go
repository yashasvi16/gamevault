package handler

import (
	"net/http"
	"github.com/yashasvi16/gamevault/internal/repository"
	"github.com/yashasvi16/gamevault/internal/ai"
	"encoding/json"
	"log/slog"
)

type AIHandler struct {
	advisor *ai.Advisor
	playerRepo repository.PlayerRepo
	matchRepo *repository.MatchRepository
}

func NewAIHandler(advisor *ai.Advisor, playerRepo repository.PlayerRepo,
matchRepo *repository.MatchRepository) *AIHandler {
	return &AIHandler{
		advisor: advisor,
		playerRepo: playerRepo,
		matchRepo: matchRepo,
	}
}

func (h *AIHandler) GetAdvice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Get the authenticated player's ID from context (set by auth middleware)
	playerIDFloat, ok := r.Context().Value("player_id").(float64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unauthorized",
		})
		return
	}

	playerID := int(playerIDFloat)

	// 2. Retrieve - Fetch player stats from database
	player, err := h.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to fetch player data",
		})
		return
	}

	// 3. Retrive - Fetch the match history from database
	matches, err := h.matchRepo.GetMatchesByPlayerID(playerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Failed to fetch match history",
        })
		return
	}

	// 4. Augment + Generate - Send data to LLM
	slog.Info("requesting AI advice", "player_id", playerID, "matches_count", len(matches))
	advice, err := h.advisor.GetAdvice(r.Context(), player, matches)
	if err != nil {
		slog.Error("AI advice failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "AI advisor is currently unavailable",
		})
		return
	}

	// 5. Return the advice
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"advice": advice,
	})
}