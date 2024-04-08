package database_sql

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	UserLevel struct {
		Level string
		Count int
	}

	Dummy struct {
		db *sql.DB
	}
)

func (d *Dummy) GoodSelectQueryContext(ctx context.Context, userID int, country string) ([]UserLevel, error) {
	const query = `SELECT level, COUNT(id) FROM users WHERE user_id = ? AND country = ? GROUP BY level ORDER BY level DESC LIMIT 15`
	rows, err := d.db.QueryContext(ctx, query, userID, country)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var ( //nolint:prealloc
		levels []UserLevel
		l      UserLevel
	)
	for rows.Next() {
		err = rows.Scan(&l.Level, &l.Count)
		if err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}

	return levels, rows.Err()
}

func (d *Dummy) GoodInsertExec(userID int, country string) error {
	const query = `INSERT INTO users (user_id, country) VALUES (?, ?)`
	_, err := d.db.Exec(query, userID, country)
	if err != nil {
		return err
	}

	return err
}

func (d *Dummy) BadInsertExec(userID int, country string) error {
	const query = `INSERT INTO users (user_id, country) VALUES (?, %s)`
	_, err := d.db.Exec(query, userID, country)
	if err != nil {
		return err
	}

	return err
}

func (d *Dummy) FalsePositiveInsertExec(userID int, country string) error {
	const query = `INSERT INTO users (user_id, country) VALUES (?, %s)`
	_, err := d.db.Exec(fmt.Sprintf(query, country), userID) //nolint:mysql
	if err != nil {
		return err
	}

	return err
}

func (d *Dummy) BadSelectQueryContext(ctx context.Context, userID int, country string) ([]UserLevel, error) {
	query := `SELECT level, COUNT(id) FROM users WHERE user_id = ? AND country = %s GROUP BY level ORDER BY level DESC LIMIT 15`
	query = query + " AND K = 4" //nolint:gocritic
	rows, err := d.db.QueryContext(ctx, query, userID, country)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var ( //nolint:prealloc
		levels []UserLevel
		l      UserLevel
	)
	for rows.Next() {
		err = rows.Scan(&l.Level, &l.Count)
		if err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}

	return levels, rows.Err()
}
