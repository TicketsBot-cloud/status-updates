// Package model contains data structures representing core domain entities.
package model

import (
	"fmt"
	"time"

	"github.com/TicketsBot-cloud/status-updates/internal/db"
)

// IncidentInfo represents the state of a Discord message for an incident
type IncidentInfo struct {
	Id            string    `json:"id" db:"id"`
	RoleId        uint64    `json:"role_id" db:"role_id"`
	MessageId     uint64    `json:"message_id" db:"message_id"`
	ThreadId      uint64    `json:"thread_id" db:"thread_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	CurrentStatus string    `json:"status" db:"status"`
}

func (i IncidentInfo) Save() error {
	_, err := db.Client.NamedExec(`INSERT INTO incidents (id, role_id, message_id, thread_id, created_at, updated_at, status)
		VALUES (:id, :role_id, :message_id, :thread_id, :created_at, :updated_at, :status)
		ON CONFLICT (id) DO UPDATE SET role_id = EXCLUDED.role_id, message_id = EXCLUDED.message_id, thread_id = EXCLUDED.thread_id,
		created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at, status = EXCLUDED.status`, i)

	if err != nil {
		fmt.Printf("Error saving incident: %v\n", err)
		return err
	}

	return nil
}
