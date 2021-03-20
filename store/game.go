package store

import (
	"context"

	log "github.com/sirupsen/logrus"

	"strconv"
	"strings"
	"time"

	"miikka.xyz/rs/config"
	"miikka.xyz/rs/graph/model"
)

// CreateGameIOT is for microcontroller. Device's client is coded to only read two players, otherwise
// schema accepts different amount of players
func (db *DB) CreateGameIOT(winner string, loser string, sport string) error {
	insertGameQuery := "insert into " + config.GameTable + " (sport_name) values($1) returning id"
	var gameID int64
	ctx := context.Background()

	// Begin transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	// Create game
	//sport = strings.ToLower(sport)
	err = db.QueryRow(ctx, insertGameQuery, sport).Scan(&gameID)
	if err != nil {
		// Must call here also or db connection drops in case of error
		tx.Rollback(ctx)
		return err
	}
	log.Info("created game with id ", gameID)

	// Create first player (winner) for a game
	winner = strings.ToLower(winner)
	playersQuery := "insert into " + config.GamePlayersTable + " (game_id, user_id, result) values($1, $2, $3) returning id"
	var playerID int64
	err = db.QueryRow(ctx, playersQuery, gameID, winner, 1).Scan(&playerID)
	if err != nil {
		// We don't want previous data to be stored if this fails
		tx.Rollback(ctx)
		return err
	}
	log.Info("added player ", winner, " to game ", gameID)

	// Create second player (loser) for a game
	loser = strings.ToLower(loser)
	err = db.QueryRow(ctx, playersQuery, gameID, loser, 0).Scan(&playerID)
	if err != nil {
		// We don't want previous data to be stored if this fails
		tx.Rollback(ctx)
		return err
	}
	log.Info("added player ", loser, " to game ", gameID)

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetUsernameByRFID return username for rf_id
func (db *DB) GetUsernameByRFID(rfid string) (*string, error) {
	query := "select user_id from " + config.UserTable + " where rf_id = $1"

	var username string
	err := db.QueryRow(context.Background(), query, rfid).Scan(&username)
	if err != nil {
		return nil, err
	}

	return &username, nil
}

// GetMostRecentGames ...
func (db *DB) GetMostRecentGames(count int) ([]*model.Game, error) {
	usersMap, err := db.GetAllUsers()
	if err != nil {
		return nil, err
	}

	if count <= 0 {
		count = 5
	}
	if count >= 50 {
		count = 50
	}

	// Each game has 2 players.
	// TODO: Later add support for multiple players per game
	count = count * 2
	query := "SELECT rs_game_players.game_id, rs_game_players.user_id, rs_game_players.result, rs_game.sport_name, rs_game.created_at FROM rs_game_players LEFT JOIN rs_game ON rs_game_players.game_id = rs_game.id where rs_game_players.game_id in (select game_id from rs_game_players) order by rs_game.id desc limit $1"
	row, err := db.Query(context.Background(), query, count)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var gamesMap = make(map[int]*model.Game)

	var games []*model.Game
	var gameID int
	var username string
	var result int
	var sport string
	var createdAt time.Time
	for row.Next() {
		row.Scan(&gameID, &username, &result, &sport, &createdAt)

		game, found := gamesMap[gameID]
		if !found {
			game = &model.Game{}
		}

		if result == config.StatusWon {
			game.Winner = usersMap[username]
		}
		if result == config.StatusLoss {
			game.Loser = usersMap[username]
		}

		game.ID = strconv.Itoa(gameID)
		game.Sport = &model.Sport{}
		game.Sport.Name = sport
		game.CreatedAt = createdAt.String()

		gamesMap[gameID] = game
		if !found {
			games = append(games, game)
		}
	}
	return games, nil
}

// Deprecated. For now adding games is allowed only with Arduino
// CreateGame creates new game
func (db *DB) CreateGame(sportName string, playersMap map[string]int) error {
	insertGameQuery := "insert into " + config.GameTable + " (sport_name) values ($1) returning id"
	var gameID int64

	ctx := context.Background()

	// Begin transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	// Create game
	sportName = strings.ToLower(sportName)
	err = db.QueryRow(ctx, insertGameQuery, sportName).Scan(&gameID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	log.Info("created game with id ", gameID)

	// Create first player (winner) for a game
	playersQuery := "insert into " + config.GamePlayersTable + " (game_id, user_id, result) values($1, $2, $3) returning id"
	var playerID int64

	for name, result := range playersMap {
		name = strings.ToLower(name)
		err = db.QueryRow(ctx, playersQuery, gameID, name, result).Scan(&playerID)
		if err != nil {
			// We don't want previous data to be stored if this fails
			defer tx.Rollback(ctx)
			return err
		}
		log.Info("added player ", name, " to game ", gameID)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
