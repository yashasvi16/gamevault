package repository

import "github.com/yashasvi16/gamevault/internal/model"

type PlayerRepo interface {
    CreatePlayer(player *model.Player) error
    GetPlayerByEmail(email string) (*model.Player, error)
    GetLeaderboard(limit, offset int) ([]model.Player, error)
	UpdatePlayerStats(playerID int, won bool) error
}
