package models

// Project defines project object
type Project struct {
	ID                 string  `json:"_id"`
	Name               string  `json:"name"`
	ColorTag           string  `json:"color_tag"`
	IsPublic           bool    `json:"is_public"`
	TrackedHours       float64 `json:"tracked_hours"`
	TrackedAmount      float64 `json:"tracked_amount"`
	ProgressPercentage float32 `json:"progress_percentage"`

	Client      *Client       `json:"client_id"`
	Workspace   *Workspace    `json:"workspace_id"`
	TeamMembers []*TeamMember `json:"-"`
	TeamGroups  []*TeamGroup  `json:"-"`
}
