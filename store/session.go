package store

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx"
	"miikka.xyz/rs/config"
)

// GetSessionToken ...
func (db *DB) GetSessionToken(username string) (*string, error) {
	const query = "select token from " + config.SessionTable + " where user_id = $1"

	var token string
	err := db.QueryRow(context.Background(), query, username).Scan(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// WriteSession ...
func (db *DB) WriteSession(username string, token string) error {
	const query = "insert into " + config.SessionTable + "(user_id, token) values($1, $2) returning id"

	var id int64
	err := db.QueryRow(context.Background(), query, username, token).Scan(&id)
	if err != nil {
		if pgerr, found := err.(pgx.PgError); found {
			// Database schema allows only uniqe values
			if pgerr.Code == "23505" {
				log.Info("Already logged in")
				return nil
			}
		}
		return err
	}

	log.Info("session inserted with id", id)

	return nil
}

// DeleteSession ...
func (db *DB) DeleteSession(token string) error {
	const query = "delete from " + config.SessionTable + " where token = $1 returning id"

	var id int64
	err := db.QueryRow(context.Background(), query, token).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

// GetUsernameFromToken validates token
func (db *DB) GetUsernameFromToken(token string) (*string, error) {
	const query = "select user_id from " + config.SessionTable + " where token = $1"

	var username string
	err := db.QueryRow(context.Background(), query, token).Scan(&username)
	if err != nil {
		return nil, err
	}

	if username == "" {
		return nil, errors.New("token was invalid")
	}

	return &username, nil
}

// SaveAudit saves event to audit trail
func (db *DB) SaveAudit(event string, subject string, object string) error {
	const query = "insert into " + config.AuditTable + "(event, subject, object) values($1, $2, $3) returning id"

	var id int64
	err := db.QueryRow(context.Background(), query, event, subject, object).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
