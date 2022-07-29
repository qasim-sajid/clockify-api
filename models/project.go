package models

// Project defines project object
type Project struct {
	ID                 string  `json:"_id"`
	Name               string  `json:"name"`
	ColorTag           string  `json:"color_tag"`
	IsPublic           bool    `json:"is_public"`
	TrackedHours       float64 `json:"tracked_hours"`
	TrackedAmount      float64 `json:"tracked_amount"`
	ProgressPercentage float64 `json:"progress_percentage"`

	Client      string   `json:"client_id"`
	Workspace   string   `json:"workspace_id"`
	TeamMembers []string `json:"team_members"`
	TeamGroups  []string `json:"team_groups"`
}
