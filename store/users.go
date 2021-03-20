package store

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"miikka.xyz/rs/config"
	"miikka.xyz/rs/graph/model"
)

// GetAllUsers returns all users in a map where username is key
func (db *DB) GetAllUsers() (map[string]*model.User, error) {
	var usersMap = make(map[string]*model.User)

	query := "select user_id, first_name, last_name from " + config.UserTable

	row, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		user := model.User{}
		row.Scan(&user.Username, &user.FirstName, &user.LastName)
		usersMap[user.Username] = &user
	}
	return usersMap, nil
}

// GetUsersWithPrefix returns all users that has prefix
func (db *DB) GetUsersWithPrefix(prefix string, pageNumber int, pageSize int) ([]*model.User, error) {
	if pageNumber < 0 {
		pageNumber = 0
	}
	if pageSize < 0 {
		pageSize = 0
	}

	var row pgx.Rows
	var err error
	//row = nil

	// Full name search
	if has, names := splitSpaceSeperatedStrings(prefix); has == true {
		query := "select user_id, first_name, last_name from " + config.UserTable + " where upper(first_name) like upper($1 || '%') and upper(last_name) like upper($2 || '%') order by first_name, last_name OFFSET $3 LIMIT $4"
		row, err = db.Query(context.Background(), query, names[0], names[1], pageNumber, pageSize)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	} else {
		query := "select user_id, first_name, last_name from " + config.UserTable + " where upper(user_id) like upper($1 || '%') or upper(first_name) like upper($1 || '%') or upper(last_name) like upper($1 || '%') order by first_name, last_name OFFSET $2 LIMIT $3"
		row, err = db.Query(context.Background(), query, prefix, pageNumber, pageSize)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	var users []*model.User
	for row.Next() {
		user := model.User{}
		row.Scan(&user.Username, &user.FirstName, &user.LastName)
		users = append(users, &user)
	}
	row.Close()
	return users, nil
}

// GetUserPassword returns crypted password for username
func (db *DB) GetUserPassword(username string) (*string, error) {
	const query = "select password from " + config.UserTable + " where user_id = $1"

	var password string
	err := db.QueryRow(context.Background(), query, username).Scan(&password)
	if err != nil {
		return nil, err
	}

	return &password, nil
}

// GetUsernameByTag ...
func (db *DB) GetUsernameByTag(tag string) (*string, error) {
	const query = "select user_id from " + config.UserTable + " where rf_id = $1"

	var username string
	err := db.QueryRow(context.Background(), query, tag).Scan(&username)
	if err != nil {
		return nil, err
	}

	return &username, nil
}

// AddRFID adds tag for player
func (db *DB) AddRFID(username string, tag string) (*string, error) {
	if tag == "" || username == "" {
		return nil, errors.New("empty value")
	}

	const query = "update " + config.UserTable + " set rf_id = $1 where user_id = $2 returning rf_id"
	var rfid string
	err := db.QueryRow(context.Background(), query, tag, username).Scan(&rfid)
	if err != nil {
		return nil, err
	}

	return &rfid, nil
}

// IsUsernameAvailable ...
func (db *DB) IsUsernameAvailable(userID string) bool {
	const query = "select count(*) as cnt from " + config.UserTable + " where user_id = $1"
	rows, err := db.Query(context.Background(), query, userID)
	if err != nil {
		log.Error(err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		rows.Scan(&count)
		return count <= 0
	}
	return true
}

// CreateUser ...
func (db *DB) CreateUser(userID string, firstName string, LastName string, email string) (int, error) {
	const query = "insert into " + config.UserTable + " (user_id, first_name, last_name, email) values($1,$2,$3,$4) returning id"
	var id int
	err := db.QueryRow(context.Background(), query, userID, firstName, LastName, email).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func splitSpaceSeperatedStrings(str string) (bool, []string) {
	if !strings.Contains(str, " ") {
		return false, nil
	}
	strs := strings.SplitN(str, " ", 2)
	if len(strs) != 2 {
		return false, nil
	}
	return true, strs
}

// GetLatestTags returns latest unregistered tags
func (db *DB) GetLatestTags(n int) ([]*model.TagAvailable, error) {
	const query = "select subject, created_at from " + config.AuditTable + " where event = 'unknown_tag' and subject not in (select rf_id from rs_user where rf_id != null) order by created_at desc;"
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	tagsMap := make(map[string]*model.TagAvailable)
	var tags []*model.TagAvailable
	for rows.Next() {
		var tag string
		var createdAt time.Time
		rows.Scan(&tag, &createdAt)

		obj := &model.TagAvailable{}
		obj.CreatedAt = createdAt.String()
		obj.ID = tag

		tagsMap[tag] = obj
		if len(tagsMap) == n {
			break
		}
	}

	for _, v := range tagsMap {
		tags = append(tags, v)
	}

	return tags, nil
}
