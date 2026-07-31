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