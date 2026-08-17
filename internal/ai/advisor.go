package ai

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/yashasvi16/gamevault/internal/model"
	"google.golang.org/api/option"
)

type Advisor struct {
	client *genai.Client
}

func NewAdvisor(apiKey string) (*Advisor, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &Advisor{client: client}, nil
}

func (a *Advisor) GetAdvice(ctx context.Context, player *model.Player, matches []model.Match) (string, error) {
	//Retrieve data from Database (player stats + match history)
	//and Augment the prompt with it, so the LLM can generate a grounded answer

	prompt := buildPrompt(player, matches)
	model := a.client.GenerativeModel("gemini-3.6-flash")
	model.SetTemperature(0.7)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gemini API call failed: %w", err)
	}

	//Extract text from response
	if len(resp.Candidates) == 0 {
		return "No advice available at this time.", nil
	}

	//Get the text content from the first candidate
	part := resp.Candidates[0].Content.Parts[0]
	return fmt.Sprintf("%v", part), nil
	
}

func buildPrompt(player *model.Player, matches []model.Match) string {
	matchSummary := ""
	for _, m := range matches {
		won := m.WinnerID != nil && *m.WinnerID == player.ID
		result := "LOST"
		if won {
			result = "WON"
		}

		var myScore, oppScore int
		if m.Player1ID == player.ID {
			myScore = m.Player1Score
			oppScore = m.Player2Score
		} else {
			myScore = m.Player2Score
			oppScore = m.Player1Score
		}

		matchSummary += fmt.Sprintf("- %s: Score %d vs %d\n", result, myScore, oppScore)
	}

	return fmt.Sprintf(`You are a gaming coach analyzing a player's performance.
	
	Player: %s
	Total Matches: %d
	Wins: %d
	Losses: %d
	Current Score: %d
	Win Rate: %0.1f%%
	
	Recent Match History:
	%s
	
	Based on this data, provide:
	1. A brief performance analysis (2-3 sentences)
	2. Their biggest strength
	3. Their biggest weakness
	4. One specific, actionable tip to improve
	
	Keep the response concise and encouraging. Use a coaching tone.`,
	player.Username, player.TotalMatches, player.WinsCount, player.LossesCount,
	player.Score,
	float64(player.WinsCount)/max(float64(player.TotalMatches), 1.0) * 100,
	matchSummary)
	
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}

	return b
}