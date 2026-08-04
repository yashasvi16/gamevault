package service

func DetermineWinner(player1ID int, player2ID int, player1Score int, player2Score int) *int {
	if player1Score > player2Score {
		return &player1ID
	} else if player2Score > player1Score {
		return &player2ID
	}
	return nil
}