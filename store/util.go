package store

import (
	"context"
	"io/ioutil"
	"strings"
)

func ExecuteSQLFile(db *DB, path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	for _, q := range strings.Split(string(file), ";") {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(ctx, q); err != nil {
			tx.Rollback(ctx)
			return err
		}
	}
	return tx.Commit(ctx)
}
