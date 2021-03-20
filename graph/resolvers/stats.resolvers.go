package resolvers

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"miikka.xyz/rs/graph/model"
)

func (r *queryResolver) GetStatsForUser(ctx context.Context, username string) (*model.StatsSummary, error) {
	stats, err := r.db.GetAllStats(username)
	if err != nil {
		log.Error(err)
		return nil, errors.New("fetching stats failed")
	}

	log.Info("stats for user ", username)

	return stats, nil
}
