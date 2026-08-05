package worker

import (
	"fmt"
	"github.com/yashasvi16/gamevault/internal/repository"
)

type StatsJob struct {
	PlayerID int
	Won bool
}

func StartStatsWorker(jobs <-chan StatsJob, repo repository.PlayerRepo) {
	for job := range jobs {
		err := repo.UpdatePlayerStats(job.PlayerID, job.Won)
		if err != nil {
			fmt.Println("Error updating stats for player", job.PlayerID, ":", err)
		} else {
			fmt.Println("Updated stats for player", job.PlayerID)
		}
	}
}