package store

import (
	"context"

	log "github.com/sirupsen/logrus"

	"strconv"
	"strings"

	pgx "github.com/jackc/pgx/v4"
	"miikka.xyz/rs/config"
	"miikka.xyz/rs/graph/model"
)

// Holds results (win, loss, tie) for each game. Key is game_id
type resultMap map[int]int

// Holds game_id array for each sport and opponent
// <sport_name>, <opponent_name>, [game_id]
type gameMap map[string]map[string][]int

// GetAllStats calculates all stats for player username.
// Sorted by sport and opponent name
func (db *DB) GetAllStats(username string) (*model.StatsSummary, error) {
	// This SQL-query combines all relevant info for played game
	// returns game_id, user_id, result, sport:
	// 4, tuommii, 1, Chess
	// 4, Prontto, 0, Chess
	// 5, ...
	// TODO: Maybe more code and simpler query ?
	query := "SELECT rs_game_players.game_id, rs_game_players.user_id, rs_game_players.result, rs_game.sport_name FROM rs_game_players LEFT JOIN rs_game ON rs_game_players.game_id = rs_game.id where rs_game_players.game_id in (select game_id from rs_game_players where user_id = $1) order by game_id asc"
	row, err := db.Query(context.Background(), query, username)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	games, results := readRows(row, username)
	stats := calcStats(db, games, results)

	return stats, nil
}

func readRows(row pgx.Rows, username string) (gameMap, resultMap) {
	var gameID int
	var userID string
	var statusInGame int
	var sportName string

	games := make(gameMap)
	results := make(resultMap)
	for row.Next() {
		row.Scan(&gameID, &userID, &statusInGame, &sportName)

		// Row contains opponent's data
		if userID != username {
			// Check that sport exists
			sportGames, found := games[sportName]
			if !found {
				temp := make(map[string][]int)
				sportGames = temp
			}

			// Check that opponent exists
			gameIDList, found := sportGames[userID]
			if !found {
				var slice []int
				gameIDList = slice
			}

			// Add data to games map
			gameIDList = append(gameIDList, gameID)
			sportGames[userID] = gameIDList
			games[sportName] = sportGames
		} else {
			// Store game result
			results[gameID] = statusInGame
		}
	}
	return games, results
}

func calcStats(db *DB, games gameMap, results resultMap) *model.StatsSummary {
	summary := model.StatsSummary{}

	// We set user info for each opponent
	users, err := db.GetAllUsers()
	if err != nil {
		log.Error("users failed", err)
	}

	for sportName, opponentMap := range games {
		sport := &model.StatsForSport{}
		sport.SportName = sportName

		for opponentName, gameIDList := range opponentMap {
			bs := calcStatsForPlayer(opponentName, gameIDList, results)
			player := &model.StatsForPlayer{}
			// TODO: Take whole user from users map. Now username is enough
			u := &model.User{}
			u.Username = opponentName
			player.User = u
			if users != nil {
				player.User = users[opponentName]
			}
			player.Stats = bs
			sport.Stats = append(sport.Stats, player)
		}
		summary.Stats = append(summary.Stats, sport)
	}
	return &summary
}

func calcStatsForPlayer(name string, gameIDList []int, results resultMap) *model.StatsEntry {
	bs := model.StatsEntry{}
	for _, v := range gameIDList {
		result := results[v]
		if result == config.StatusWon {
			bs.Wins++
		} else if result == config.StatusTie {
			bs.Ties++
		} else if result == config.StatusLoss {
			bs.Losses++
		}
	}
	return &bs
}

// GetStats returns stats for username sorted by sport names
func (db *DB) GetStats(username string) (map[string]*model.StatsEntry, error) {
	query := "SELECT rs_game_players.result, rs_game.sport_name FROM rs_game_players left join rs_game on rs_game_players.game_id = rs_game.id where rs_game_players.user_id = $1"

	row, err := db.Query(context.Background(), query, username)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	statsBySport := make(map[string]*model.StatsEntry)
	var statusInGame int
	var sportName string
	for row.Next() {
		row.Scan(&statusInGame, &sportName)
		// TODO: Refactor
		obj, hasKey := statsBySport[sportName]
		if !hasKey {
			// Insert new object to map
			bs := &model.StatsEntry{}
			if statusInGame == 0 {
				bs.Losses++
			} else if statusInGame == 1 {
				bs.Wins++
			} else if statusInGame == 2 {
				bs.Ties++
			}
			statsBySport[sportName] = bs
		} else {
			// Modify existing object
			if statusInGame == 0 {
				obj.Losses++
			} else if statusInGame == 1 {
				obj.Wins++
			} else if statusInGame == 2 {
				obj.Ties++
			}
		}
	}

	return statsBySport, nil
}

// GetBasicStats ...
func (db *DB) GetBasicStats(username string) (*model.StatsEntry, error) {
	stats := &model.StatsEntry{}

	// TODO: alias, rename
	query := "select game_id, result from " + config.GamePlayersTable + " where user_id = $1"

	row, err := db.Query(context.Background(), query, username)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var statusInGame int
	var gameID int

	for row.Next() {
		row.Scan(&gameID, &statusInGame)
		if statusInGame == config.StatusLoss {
			stats.Losses++
		} else if statusInGame == config.StatusWon {
			stats.Wins++
		} else if statusInGame == config.StatusTie {
			stats.Ties++
		}
	}
	return stats, nil
}

// NOT IN USE. Might be used later
func (db *DB) getSportNames(idMap map[int]int) (map[int]string, error) {
	inQuery := "("
	for k := range idMap {
		inQuery += strconv.Itoa(k) + ","
	}
	inQuery = strings.TrimSuffix(inQuery, ",")
	inQuery += ")"

	query := "select id, sport_name from rs_game where id in " + inQuery
	row, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var gameID int
	var sportName string

	// game_id, sportName
	sportNameMap := make(map[int]string)
	for row.Next() {
		row.Scan(&gameID, &sportName)
		sportNameMap[gameID] = sportName
	}
	return sportNameMap, nil
}
