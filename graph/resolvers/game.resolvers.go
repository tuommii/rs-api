package resolvers

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"strconv"

	"miikka.xyz/rs/graph/model"
)

const maxPlayersPerGame int = 16

func (r *mutationResolver) CreateGame(ctx context.Context, input model.NewGame) (*string, error) {
	// Validate amount of players
	playersCount := len(input.Players)
	if playersCount == 0 || playersCount > maxPlayersPerGame {
		return nil, errors.New("amount (" + strconv.Itoa(playersCount) + ") of players not allowed")
	}

	if input.SportName == "" {
		return nil, errors.New("empty sport name")
	}

	// Make sure only unique players are added
	players := make(map[string]int)
	for _, player := range input.Players {
		players[player.Username] = player.Result
	}

	err := r.db.CreateGame(input.SportName, players)
	if err != nil {
		log.Error(err)
		return nil, errors.New("create game failed")
	}

	ret := "game created"
	return &ret, nil
}

func (r *queryResolver) GetMostRecentGames(ctx context.Context, count int) ([]*model.Game, error) {
	games, err := r.db.GetMostRecentGames(count)
	if err != nil {
		log.Error(err)
		return nil, errors.New("most recent games failed")
	}
	log.Info("returning ", len(games), " most recent games")
	return games, nil
}
