package store

import (
	"context"

	"miikka.xyz/rs/config"
)

// SetConfig adds new setting to database
func (db *DB) SetConfig(key string, value string, valueType string) error {
	insertQuery := "insert into " + config.ConfigTable + " (key, value, value_type) values ($1, $2) returning id"
	var id int64
	err := db.QueryRow(context.Background(), insertQuery, key, value, valueType).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

// GetConfig ...
func GetConfig(key string) error {
	return nil
}
