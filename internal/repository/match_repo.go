package repository

import (
	"database/sql"
	"github.com/yashasvi16/gamevault/internal/model"
	"fmt"
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

func (r *MatchRepository) RecordMatchWithStats(match *model.Match) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO matches (player1_id, player2_id, winner, 
	player1_score, player2_score) VALUES ($1, $2, $3, $4, $5)`,
		match.Player1ID, match.Player2ID, match.WinnerID,
		match.Player1Score, match.Player2Score)
	if err != nil {
		return fmt.Errorf("failed to insert match: %w", err)
	}

	p1Won := match.WinnerID != nil && *match.WinnerID == match.Player1ID
	_, err = tx.Exec(`UPDATE players SET total_matches = total_matches + 1,
		wins_count = CASE WHEN $2 THEN wins_count + 1 ELSE wins_count END,
		losses_count = CASE WHEN $2 THEN losses_count ELSE losses_count + 1 END,
		score = CASE WHEN $2 THEN score + 10 ELSE score - 5 END
		WHERE id = $1`, match.Player1ID, p1Won)
	if err != nil {
		return fmt.Errorf("failed to update player1 stats: %w", err)
	}

	p2Won := match.WinnerID != nil && *match.WinnerID == match.Player2ID
	_, err = tx.Exec(`UPDATE players SET total_matches = total_matches + 1,
		wins_count = CASE WHEN $2 THEN wins_count + 1 ELSE wins_count END,
		losses_count = CASE WHEN $2 THEN losses_count ELSE losses_count + 1 END,
		score = CASE WHEN $2 THEN score + 10 ELSE score - 5 END
		WHERE id = $1`, match.Player2ID, p2Won)
	if err != nil {
		return fmt.Errorf("failed to update player2 stats: %w", err)
	}

	return tx.Commit()
}