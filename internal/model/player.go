package model

import "time"

type Player struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	AvatarURL string `json:"avatar_url"`
	TotalMatches int `json:"total_matches"`
	WinsCount int `json:"wins_count"`
	LossesCount int `json:"losses_count"`
	Score int `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}