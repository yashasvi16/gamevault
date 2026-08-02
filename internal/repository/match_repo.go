package repository

import (
	"database/sql"
	"github.com/yashasvi16/gamevault/internal/model"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository {
		db: db,
	}
}

func (r *MatchRepository) CreateMatch(match *model.Match) error {
	query := `INSERT INTO matches (player1_id, player2_id, winner, 
	player1_score, player2_score) VALUES ($1, $2, $3, $4, $5) RETURNING id,
	created_at`

	err := r.db.QueryRow(query, match.Player1ID, match.Player2ID, match.WinnerID,
	match.Player1Score, match.Player2Score).Scan(&match.ID, &match.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}