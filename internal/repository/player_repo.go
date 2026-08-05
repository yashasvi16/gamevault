package repository

import (
	"database/sql"
	"github.com/yashasvi16/gamevault/internal/model"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

func (r *PlayerRepository) CreatePlayer(player *model.Player) error {
	query := `INSERT INTO players(username, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id, created_at`

	err := r.db.QueryRow(query, player.Username, player.Email, player.PasswordHash).Scan(&player.ID, &player.CreatedAt)
	return err
}

func (r *PlayerRepository) GetPlayerByEmail(email string) (*model.Player, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(avatar_url, ''),
	total_matches, wins_count, losses_count, score, created_at FROM players
	WHERE email = $1`

	var player model.Player
	err := r.db.QueryRow(query, email).Scan(&player.ID, &player.Username, &player.Email, 
	&player.PasswordHash, &player.AvatarURL, &player.TotalMatches, 
	&player.WinsCount, &player.LossesCount, &player.Score, &player.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (r *PlayerRepository) GetLeaderboard(limit int, offset int) ([]model.Player, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(avatar_url, ''),
	total_matches, wins_count, losses_count, score, created_at FROM players ORDER BY score DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	

	players := []model.Player{}
	for rows.Next() {
		var player model.Player
		
		err = rows.Scan(&player.ID, &player.Username, &player.Email, 
		&player.PasswordHash, &player.AvatarURL, &player.TotalMatches, 
		&player.WinsCount, &player.LossesCount, &player.Score, &player.CreatedAt)
		
		if err != nil {
			return nil, err
		}

		players = append(players, player)

	}

	return players, nil
}

func (r *PlayerRepository) UpdatePlayerStats(playerID int, won bool) error {
	query := `UPDATE players 
	SET total_matches = total_matches + 1,
		wins_count = CASE 
			WHEN $2 THEN wins_count + 1 
			ELSE wins_count
		END,
		losses_count = CASE
			WHEN $2 THEN losses_count
			ELSE losses_count + 1
		END,
		score = CASE
			WHEN $2 THEN score + 10
			ELSE score - 5
		END
	WHERE id = $1`

	_, err := r.db.Exec(query, playerID, won)
	
	if err !=  nil {
		return err
	}

	return nil
}