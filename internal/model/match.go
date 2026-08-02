package model

import "time"

type Match struct {
	ID           int       `json:"id"`
	Player1ID    int       `json:"player1_id"`
	Player2ID    int       `json:"player2_id"`
	WinnerID     *int      `json:"winner_id"` // Pointer to int allows NULL when match is ongoing or draw
	Player1Score int       `json:"player1_score"`
	Player2Score int       `json:"player2_score"`
	CreatedAt    time.Time `json:"created_at"`
}