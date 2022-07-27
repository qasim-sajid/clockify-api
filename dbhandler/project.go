package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddProject(project *models.Project) (*models.Project, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("Unable to generate _ID")
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

func (db *dbClient) AddProjectTeamMembers(projectID string, teamMembers []*models.TeamMember) error {
	if teamMembers == nil {
		return nil
	}

	for _, tm := range teamMembers {
		valuesMap := make(map[string]interface{})
		valuesMap["project_id"] = projectID
		valuesMap["team_member_id"] = tm.ID

		//Check if value already exists
		_, err := db.GetTeamMemberForProject(projectID, tm.ID)
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

func (db *dbClient) GetTeamMemberForProject(projectID, teamMemberID string) (*models.TeamMember, error) {
	teamMembers, err := db.GetProjectTeamMembers(projectID)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMemberForProject: %v", err)
	}

	for _, tm := range teamMembers {
		if tm.ID == teamMemberID {
			return tm, nil
		}
	}

	return nil, fmt.Errorf("GetTeamMemberForProject: %v", errors.New("TeamMember with given ID not found!"))
}

func (db *dbClient) AddProjectTeamGroups(projectID string, teamGroups []*models.TeamGroup) error {
	if teamGroups == nil {
		return nil
	}

	for _, tg := range teamGroups {
		valuesMap := make(map[string]interface{})
		valuesMap["project_id"] = projectID
		valuesMap["team_group_id"] = tg.ID

		//Check if value already exists
		_, err := db.GetTeamGroupForProject(projectID, tg.ID)
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

func (db *dbClient) GetTeamGroupForProject(projectID, teamGroupID string) (*models.TeamGroup, error) {
	teamGroups, err := db.GetProjectTeamGroups(projectID)
	if err != nil {
		return nil, fmt.Errorf("GetTeamGroupForProject: %v", err)
	}

	for _, tg := range teamGroups {
		if tg.ID == teamGroupID {
			return tg, nil
		}
	}

	return nil, fmt.Errorf("GetTeamGroupForProject: %v", errors.New("TeamGroup with given ID not found!"))
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

	rows, err := db.GetValuesFromCompositeTable(PROJECT_TEAM_MEMBER, searchParams)
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

	rows, err := db.GetValuesFromCompositeTable(PROJECT_TEAM_GROUP, searchParams)
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

	if v, ok := updates["project_team_members"]; ok {
		teamMembers := v.([]*models.TeamMember)
		err = db.UpdateProjectTeamMembers(projectID, teamMembers)
		if err != nil {
			return nil, fmt.Errorf("UpdateProject: %v", err)
		}
	}

	if v, ok := updates["project_team_groups"]; ok {
		teamGroups := v.([]*models.TeamGroup)
		err = db.UpdateProjectTeamGroups(projectID, teamGroups)
		if err != nil {
			return nil, fmt.Errorf("UpdateProject: %v", err)
		}
	}

	return project, nil
}

func (db *dbClient) UpdateProjectTeamMembers(projectID string, teamMembers []*models.TeamMember) error {
	deleteParams := make(map[string]interface{})
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

func (db *dbClient) UpdateProjectTeamGroups(projectID string, teamGroups []*models.TeamGroup) error {
	deleteParams := make(map[string]interface{})
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
