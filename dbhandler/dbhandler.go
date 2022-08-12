package dbhandler

import (
	"github.com/qasim-sajid/clockify-api/models"
)

// DbHandler specifies DB context
type DbHandler interface {
	SetupDB()
	CloseDB()

	AddClient(*models.Client) (*models.Client, int, error)
	GetAllClients() ([]*models.Client, error)
	GetClientsWithFilters(searchParams map[string]interface{}) ([]*models.Client, error)
	GetClient(clientID string) (*models.Client, error)
	UpdateClient(clientID string, updates map[string]interface{}) (*models.Client, error)
	DeleteClient(clientID string) error

	AddProject(*models.Project) (*models.Project, int, error)
	GetAllProjects() ([]*models.Project, error)
	GetProjectsWithFilters(searchParams map[string]interface{}) ([]*models.Project, error)
	GetProject(projectID string) (*models.Project, error)
	UpdateProject(projectID string, updates map[string]interface{}) (*models.Project, error)
	DeleteProject(projectID string) error

	AddTag(*models.Tag) (*models.Tag, int, error)
	GetAllTags() ([]*models.Tag, error)
	GetTagsWithFilters(searchParams map[string]interface{}) ([]*models.Tag, error)
	GetTag(tagID string) (*models.Tag, error)
	UpdateTag(tagID string, updates map[string]interface{}) (*models.Tag, error)
	DeleteTag(tagID string) error

	AddTask(*models.Task) (*models.Task, int, error)
	GetAllTasks() ([]*models.Task, error)
	GetTasksWithFilters(searchParams map[string]interface{}) ([]*models.Task, error)
	GetTask(taskID string) (*models.Task, error)
	UpdateTask(taskID string, updates map[string]interface{}) (*models.Task, error)
	DeleteTask(taskID string) error

	AddTeamGroup(*models.TeamGroup) (*models.TeamGroup, int, error)
	GetAllTeamGroups() ([]*models.TeamGroup, error)
	GetTeamGroupsWithFilters(searchParams map[string]interface{}) ([]*models.TeamGroup, error)
	GetTeamGroup(teamGroupID string) (*models.TeamGroup, error)
	UpdateTeamGroup(teamGroupID string, updates map[string]interface{}) (*models.TeamGroup, error)
	DeleteTeamGroup(teamGroupID string) error

	AddTeamMember(*models.TeamMember) (*models.TeamMember, int, error)
	AddTeamMemberTeamGroups(string, []string) error
	GetAllTeamMembers() ([]*models.TeamMember, error)
	GetTeamMembersWithFilters(searchParams map[string]interface{}) ([]*models.TeamMember, error)
	GetTeamMember(teamMemberID string) (*models.TeamMember, error)
	UpdateTeamMember(teamMemberID string, updates map[string]interface{}) (*models.TeamMember, error)
	DeleteTeamMember(teamMemberID string) error

	AddTeamRole(*models.TeamRole) (*models.TeamRole, int, error)
	GetAllTeamRoles() ([]*models.TeamRole, error)
	GetTeamRolesWithFilters(searchParams map[string]interface{}) ([]*models.TeamRole, error)
	GetTeamRole(teamRoleID string) (*models.TeamRole, error)
	UpdateTeamRole(teamRoleID string, updates map[string]interface{}) (*models.TeamRole, error)
	DeleteTeamRole(teamRoleID string) error

	AddUser(*models.User) (*models.User, int, error)
	GetAllUsers() ([]*models.User, error)
	GetUsersWithFilters(searchParams map[string]interface{}) ([]*models.User, error)
	GetUser(userID string) (*models.User, error)
	GetUserWithIdentity(userID string) (*models.User, error)
	UpdateUser(userID string, updates map[string]interface{}) (*models.User, error)
	DeleteUser(userID string) error
	CheckUserLogin(string, string) (*models.User, error)

	AddWorkspace(*models.Workspace) (*models.Workspace, int, error)
	GetAllWorkspaces() ([]*models.Workspace, error)
	GetWorkspacesWithFilters(searchParams map[string]interface{}) ([]*models.Workspace, error)
	GetWorkspace(workspaceID string) (*models.Workspace, error)
	UpdateWorkspace(workspaceID string, updates map[string]interface{}) (*models.Workspace, error)
	DeleteWorkspace(workspaceID string) error
}
