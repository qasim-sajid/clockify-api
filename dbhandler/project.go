package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddProject(project *models.Project) (*models.Project, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("unable to generate _id")
	}
	project.ID = fmt.Sprintf("p_%v", id)

	insertQuery, err := db.GetInsertQuery(*project)
	if err != nil {
		return nil, -1, fmt.Errorf("AddProject: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddProject: %v", err)
	}

	err = db.AddProjectTeamMembers(project.ID, project.TeamMembers)
	if err != nil {
		return nil, -1, fmt.Errorf("AddProject: %v", err)
	}

	err = db.AddProjectTeamGroups(project.ID, project.TeamGroups)
	if err != nil {
		return nil, -1, fmt.Errorf("AddProject: %v", err)
	}

	return project, http.StatusOK, nil
}

func (db *dbClient) AddProjectTeamMembers(projectID string, teamMembers []string) error {
	if teamMembers == nil || teamMembers[0] == "" {
		return nil
	}

	for _, tm := range teamMembers {
		valuesMap := make(map[string]interface{})
		valuesMap["project_id"] = projectID
		valuesMap["team_member_id"] = tm

		//Check if value already exists
		_, err := db.GetTeamMemberForProject(projectID, tm)
		if err != nil {
			//If value doesn't exist then insert it
			_, err := db.AddValueInCompositeTable(PROJECT_TEAM_MEMBER, valuesMap)
			if err != nil {
				return fmt.Errorf("AddProjectTeamMembers: %v", err)
			}
		}
	}

	return nil
}

func (db *dbClient) GetTeamMemberForProject(projectID, teamMemberID string) (string, error) {
	teamMembers, err := db.GetProjectTeamMembers(projectID)
	if err != nil {
		return "", fmt.Errorf("GetTeamMemberForProject: %v", err)
	}

	for _, tm := range teamMembers {
		if tm == teamMemberID {
			return tm, nil
		}
	}

	return "", fmt.Errorf("GetTeamMemberForProject: %v", errors.New("team member with given _id not found"))
}

func (db *dbClient) AddProjectTeamGroups(projectID string, teamGroups []string) error {
	if teamGroups == nil || teamGroups[0] == "" {
		return nil
	}

	for _, tg := range teamGroups {
		valuesMap := make(map[string]interface{})
		valuesMap["project_id"] = projectID
		valuesMap["team_group_id"] = tg

		//Check if value already exists
		_, err := db.GetTeamGroupForProject(projectID, tg)
		if err != nil {
			//If value doesn't exist then insert it
			_, err := db.AddValueInCompositeTable(PROJECT_TEAM_GROUP, valuesMap)
			if err != nil {
				return fmt.Errorf("AddProjectTeamGroups: %v", err)
			}
		}
	}

	return nil
}

func (db *dbClient) GetTeamGroupForProject(projectID, teamGroupID string) (string, error) {
	teamGroups, err := db.GetProjectTeamGroups(projectID)
	if err != nil {
		return "", fmt.Errorf("GetTeamGroupForProject: %v", err)
	}

	for _, tg := range teamGroups {
		if tg == teamGroupID {
			return tg, nil
		}
	}

	return "", fmt.Errorf("GetTeamGroupForProject: %v", errors.New("team group with given id not found"))
}

func (db *dbClient) GetAllProjects() ([]*models.Project, error) {
	projects, err := db.GetProjectsWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllProjects: %v", err)
	}

	return projects, nil
}

func (db *dbClient) GetProject(projectID string) (*models.Project, error) {
	selectParams := make(map[string]interface{})

	selectParams["_id"] = projectID

	projects, err := db.GetProjectsWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetProject: %v", err)
	}

	var project *models.Project
	if projects == nil || len(projects) <= 0 {
		return nil, fmt.Errorf("GetProject: %v", errors.New("project with given if not found"))
	} else {
		project = projects[0]
	}

	return project, nil
}

func (db *dbClient) GetProjectsWithFilters(searchParams map[string]interface{}) ([]*models.Project, error) {
	p := models.Project{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetProjectsWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetProjectsWithFilters: %v", err)
	}

	projects, err := db.GetProjectsFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetProjectsWithFilters: %v", err)
	}

	return projects, nil
}

func (db *dbClient) GetProjectsFromRows(rows *sql.Rows) ([]*models.Project, error) {
	projects := make([]*models.Project, 0)
	for rows.Next() {
		p := models.Project{}

		var clientID sql.NullString
		var workspaceID sql.NullString

		err := rows.Scan(&p.ID, &p.Name, &p.ColorTag, &p.IsPublic, &p.TrackedHours, &p.TrackedAmount, &p.ProgressPercentage, &clientID,
			&workspaceID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectsFromRows: %v", err)
		}

		if clientID.Valid {
			p.Client = clientID.String
		}

		if workspaceID.Valid {
			p.Workspace = workspaceID.String
		}

		p.TeamMembers, err = db.GetProjectTeamMembers(p.ID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectsFromRows: %v", err)
		}

		p.TeamGroups, err = db.GetProjectTeamGroups(p.ID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectsFromRows: %v", err)
		}

		projects = append(projects, &p)
	}

	return projects, nil
}

func (db *dbClient) GetProjectTeamMembers(projectID string) ([]string, error) {
	searchParams := make(map[string]interface{})
	searchParams["project_id"] = projectID

	rows, err := db.GetValuesFromCompositeTable(PROJECT_TEAM_MEMBER, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetProjectTeamMembers: %v", err)
	}

	teamMembers := make([]string, 0)
	for rows.Next() {
		teamMemberID := ""

		err := rows.Scan(&projectID, &teamMemberID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectTeamMembers: %v", err)
		}

		teamMembers = append(teamMembers, teamMemberID)
	}

	return teamMembers, nil
}

func (db *dbClient) GetProjectTeamGroups(projectID string) ([]string, error) {
	searchParams := make(map[string]interface{})
	searchParams["project_id"] = projectID

	rows, err := db.GetValuesFromCompositeTable(PROJECT_TEAM_GROUP, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetProjectTeamGroups: %v", err)
	}

	teamGroups := make([]string, 0)
	for rows.Next() {
		teamGroupID := ""

		err := rows.Scan(&projectID, &teamGroupID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectTeamGroups: %v", err)
		}

		teamGroups = append(teamGroups, teamGroupID)
	}

	return teamGroups, nil
}

func (db *dbClient) UpdateProject(projectID string, updates map[string]interface{}) (*models.Project, error) {
	if v, ok := updates["team_members"]; ok {
		teamMembers := strings.Split(v.(string), ",")
		if len(teamMembers) > 0 && teamMembers[0] != "" {
			err := db.UpdateProjectTeamMembers(projectID, teamMembers)
			if err != nil {
				return nil, fmt.Errorf("UpdateProject: %v", err)
			}
		}
		delete(updates, "team_members")
	}

	if v, ok := updates["team_groups"]; ok {
		teamGroups := strings.Split(v.(string), ",")
		if len(teamGroups) > 0 && teamGroups[0] != "" {
			err := db.UpdateProjectTeamGroups(projectID, teamGroups)
			if err != nil {
				return nil, fmt.Errorf("UpdateProject: %v", err)
			}
		}
		delete(updates, "team_groups")
	}

	updateQuery, err := db.GetUpdateQueryForStruct(models.Project{}, projectID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateProject: %v", err)
	}

	if len(updates) > 0 {
		_, err = db.RunUpdateQuery(updateQuery)
		if err != nil {
			return nil, fmt.Errorf("UpdateProject: %v", err)
		}
	}

	project, err := db.GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("UpdateProject: %v", err)
	}

	return project, nil
}

func (db *dbClient) UpdateProjectTeamMembers(projectID string, teamMembers []string) error {
	deleteParams := make(map[string]interface{})
	deleteParams["team_member_id"] = projectID
	_, err := db.DeleteValuesFromCompositeTable(PROJECT_TEAM_MEMBER, deleteParams)
	if err != nil {
		return fmt.Errorf("UpdateProjectTeamMembers: %v", err)
	}

	err = db.AddProjectTeamMembers(projectID, teamMembers)
	if err != nil {
		return fmt.Errorf("UpdateProjectTeamMembers: %v", err)
	}

	return nil
}

func (db *dbClient) UpdateProjectTeamGroups(projectID string, teamGroups []string) error {
	deleteParams := make(map[string]interface{})
	deleteParams["team_group_id"] = projectID
	_, err := db.DeleteValuesFromCompositeTable(PROJECT_TEAM_GROUP, deleteParams)
	if err != nil {
		return fmt.Errorf("UpdateProjectTeamGroups: %v", err)
	}

	err = db.AddProjectTeamGroups(projectID, teamGroups)
	if err != nil {
		return fmt.Errorf("UpdateProjectTeamGroups: %v", err)
	}

	return nil
}

func (db *dbClient) DeleteProject(projectID string) error {
	deleteParams := make(map[string]interface{})

	deleteParams["_id"] = projectID

	deleteQuery, err := db.GetDeleteQueryForStruct(models.Project{}, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteProject: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteProject: %v", err)
	}

	deleteParamsForColumns := make(map[string]interface{})
	deleteParamsForColumns["project_id"] = projectID

	_, err = db.DeleteValuesFromCompositeTable(PROJECT_TEAM_GROUP, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteTeamGroupsForProject: %v", err)
	}

	_, err = db.DeleteValuesFromCompositeTable(PROJECT_TEAM_MEMBER, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteTeamMembersForProject: %v", err)
	}

	return nil
}

func (db *dbClient) AddValueInCompositeTable(tableName string, valuesMap map[string]interface{}) (sql.Result, error) {
	insertQuery, err := db.GetInsertQueryForCompositeTable(tableName, valuesMap)
	if err != nil {
		return nil, fmt.Errorf("AddValuesInCompositeTable: %v", err)
	}

	result, err := db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, fmt.Errorf("AddValuesInCompositeTable: %v", err)
	}

	return result, nil
}

func (db *dbClient) GetValuesFromCompositeTable(tableName string, searchParams map[string]interface{}) (*sql.Rows, error) {
	selectQuery, err := db.GetSelectQueryForCompositeTable(tableName, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetValuesFromCompositeTable: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetValuesFromCompositeTable: %v", err)
	}

	return rows, nil
}

func (db *dbClient) DeleteValuesFromCompositeTable(tableName string, deleteParams map[string]interface{}) (sql.Result, error) {
	deleteQuery, err := db.GetDeleteQueryForCompositeTable(tableName, deleteParams)
	if err != nil {
		return nil, fmt.Errorf("DeleteValuesInCompositeTable: %v", err)
	}

	result, err := db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return nil, fmt.Errorf("DeleteValuesInCompositeTable: %v", err)
	}

	return result, nil
}
