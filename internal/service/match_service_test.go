package service

import (
	"testing"
)

func TestDetermineWinner(t *testing.T) {
	player1ID := 1
	player2ID := 2

	tests := []struct {
		test_name string
		p1_score int
		p2_score int
		expected_winner *int
	}{
		{"player 1 wins", 10, 5, &player1ID},
		{"player 2 wins", 3, 7, &player2ID},
		{"draw", 5, 5, nil},
	}

	for _,tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			result := DetermineWinner(player1ID, player2ID, tt.p1_score, tt.p2_score)
			if (result == nil) != (tt.expected_winner == nil) {
				t.Errorf("DetermineWinner(%d, %d) = %d, want %d",
				tt.p1_score, tt.p2_score, result, tt.expected_winner)
			} else if (result != nil) && *result != *tt.expected_winner {
				t.Errorf("DetermineWinner(%d, %d) = %d, want %d",
				tt.p1_score, tt.p2_score, result, tt.expected_winner)
			}
		})
	}
}