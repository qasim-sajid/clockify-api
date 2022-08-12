package models

// TeamGroup defines team_group object
type TeamGroup struct {
	ID   string `json:"_id"`
	Name string `json:"name"`

	Workspace   string   `json:"workspace_id"`
	TeamMembers []string `json:"team_members"`
}
