package model

// Metadata represents the structure of metadata in the incident
type Metadata struct {
	JIRA JIRAData `json:"jira"`
}

// JIRAData represents the structure of JIRA-related metadata
type JIRAData struct {
	IssueID string `json:"issue_id"`
}
