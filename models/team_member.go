package models

// TeamMember defines team_member object
type TeamMember struct {
	ID           string  `json:"_id"`
	BillableRate float64 `json:"billable_rate"`

	Workspace  string   `json:"workspace_id"`
	User       string   `json:"user_email"`
	TeamRole   string   `json:"team_role_id"`
	TeamGroups []string `json:"team_groups"`
}
