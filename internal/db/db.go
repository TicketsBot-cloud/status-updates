// Package db provides database access and utility functions for the application.
package db

import (
	"github.com/TicketsBot-cloud/status-updates/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Client *sqlx.DB

func InitDB() error {
	var err error
	Client, err = sqlx.Connect("postgres", config.Conf.DatabaseUri)
	if err != nil {
		return err
	}

	// Create tables if they don't exist
	schema := `
	CREATE TABLE IF NOT EXISTS incidents (
		id TEXT PRIMARY KEY,
		role_id BIGINT NOT NULL,
		message_id BIGINT NOT NULL,
		thread_id BIGINT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		status TEXT NOT NULL
	);
	`
	_, err = Client.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}
