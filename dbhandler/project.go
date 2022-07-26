package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddProject(project *models.Project) (*models.Project, int, error) {
	insertQuery, err := db.GetInsertQueryForStruct(project)
	if err != nil {
		return nil, -1, fmt.Errorf("AddProject: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddProject: %v", err)
	}

	return project, http.StatusOK, nil
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
	p := models.Project{}
	v := reflect.ValueOf(p)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return nil, fmt.Errorf("GetProject: %v", err)
	}
	selectParams[columnName] = projectID

	projects, err := db.GetProjectsWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetProject: %v", err)
	}

	var project *models.Project
	if projects == nil || len(projects) <= 0 {
		return nil, fmt.Errorf("GetProject: %v", errors.New("Project with given ID not found!"))
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
			p.Client, err = db.GetClient(clientID.String)
			if err != nil {
				return nil, fmt.Errorf("GetProjectsFromRows: %v", err)
			}
		}

		if workspaceID.Valid {
			p.Workspace, err = db.GetWorkspace(workspaceID.String)
			if err != nil {
				return nil, fmt.Errorf("GetProjectsFromRows: %v", err)
			}
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

func (db *dbClient) GetProjectTeamMembers(projectID string) ([]*models.TeamMember, error) {
	searchParams := make(map[string]interface{})
	searchParams["project_id"] = projectID

	selectQuery, err := db.GetSelectQueryForCompositeTable(PROJECT_TEAM_MEMBER, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetProjectTeamMembers: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetProjectTeamMembers: %v", err)
	}

	teamMembers := make([]*models.TeamMember, 0)
	for rows.Next() {
		teamMemberID := ""

		err := rows.Scan(&projectID, &teamMemberID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectTeamMembers: %v", err)
		}

		tm, err := db.GetTeamMember(teamMemberID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectTeamMembers: %v", err)
		}

		teamMembers = append(teamMembers, tm)
	}

	return teamMembers, nil
}

func (db *dbClient) GetProjectTeamGroups(projectID string) ([]*models.TeamGroup, error) {
	searchParams := make(map[string]interface{})
	searchParams["project_id"] = projectID

	selectQuery, err := db.GetSelectQueryForCompositeTable(PROJECT_TEAM_GROUP, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetProjectTeamGroups: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetProjectTeamGroups: %v", err)
	}

	teamGroups := make([]*models.TeamGroup, 0)
	for rows.Next() {
		teamGroupID := ""

		err := rows.Scan(&projectID, &teamGroupID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectTeamGroups: %v", err)
		}

		tg, err := db.GetTeamGroup(teamGroupID)
		if err != nil {
			return nil, fmt.Errorf("GetProjectTeamGroups: %v", err)
		}

		teamGroups = append(teamGroups, tg)
	}

	return teamGroups, nil
}

func (db *dbClient) UpdateProject(projectID string, updates map[string]interface{}) (*models.Project, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.Project{}, projectID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateProject: %v", err)
	}

	_, err = db.RunUpdateQuery(updateQuery)
	if err != nil {
		return nil, fmt.Errorf("UpdateProject: %v", err)
	}

	project, err := db.GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("UpdateProject: %v", err)
	}

	return project, nil
}

func (db *dbClient) DeleteProject(projectID string) error {
	deleteParams := make(map[string]interface{})
	c := models.Project{}
	v := reflect.ValueOf(c)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return fmt.Errorf("DeleteProject: %v", err)
	}

	deleteParams[columnName] = projectID

	deleteQuery, err := db.GetDeleteQueryForStruct(c, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteProject: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteProject: %v", err)
	}

	return nil
}
