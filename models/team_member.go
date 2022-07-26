package models

// TeamMember defines team_member object
type TeamMember struct {
	ID           string  `json:"_id"`
	BillableRate float32 `json:"billable_rate"`

	Workspace *Workspace `json:"workspace_id"`
	User      *User      `json:"user_email"`
	TeamRole  *TeamRole  `json:"team_role_id"`
}
