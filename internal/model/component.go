package model

import "time"

// Component represents the structure of a component
type Component struct {
	ID                 string    `json:"id"`
	PageID             string    `json:"page_id"`
	GroupID            string    `json:"group_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Group              bool      `json:"group"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Position           int       `json:"position"`
	Status             string    `json:"status"`
	Showcase           bool      `json:"showcase"`
	OnlyShowIfDegraded bool      `json:"only_show_if_degraded"`
	AutomationEmail    string    `json:"automation_email"`
	StartDate          string    `json:"start_date"` // Keeping as string to preserve date format
}

// AffectedComponent represents the structure of an affected component
type AffectedComponent struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
}
