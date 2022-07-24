package models

// TeamGroup defines team_group object
type TeamGroup struct {
	ID   string `json:"_id"`
	Name string `json:"name"`

	Workspace   *Workspace    `json:"workspace_id"`
	TeamMembers *[]TeamMember `json:"-"`
}
