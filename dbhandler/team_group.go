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

func (db *dbClient) AddTeamGroup(teamGroup *models.TeamGroup) (*models.TeamGroup, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("Unable to generate _ID")
	}
	teamGroup.ID = fmt.Sprintf("tg_%v", id)

	insertQuery, err := db.GetInsertQuery(*teamGroup)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamGroup: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamGroup: %v", err)
	}

	err = db.AddTeamGroupTeamMembers(teamGroup.ID, teamGroup.TeamMembers)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamGroup: %v", err)
	}

	return teamGroup, http.StatusOK, nil
}

func (db *dbClient) AddTeamGroupTeamMembers(teamGroupID string, teamMembers []*models.TeamMember) error {
	if teamMembers == nil {
		return nil
	}

	for _, tm := range teamMembers {
		valuesMap := make(map[string]interface{})
		valuesMap["team_group_id"] = teamGroupID
		valuesMap["team_member_id"] = tm.ID

		//Check if value already exists
		_, err := db.GetTeamMemberForTeamGroup(teamGroupID, tm.ID)
		if err != nil {
			//If value doesn't exist then insert it
			insertQuery, err := db.GetInsertQueryForCompositeTable(TEAM_GROUP_TEAM_MEMBER, valuesMap)
			if err != nil {
				return fmt.Errorf("AddTeamGroupTeamMembers: %v", err)
			}

			_, err = db.RunInsertQuery(insertQuery)
			if err != nil {
				return fmt.Errorf("AddTeamGroupTeamMembers: %v", err)
			}
		}
	}

	return nil
}

func (db *dbClient) GetTeamMemberForTeamGroup(teamGroupID, teamMemberID string) (*models.TeamMember, error) {
	teamMembers, err := db.GetTeamGroupTeamMembers(teamGroupID)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMemberForTeamGroup: %v", err)
	}

	for _, tm := range teamMembers {
		if tm.ID == teamMemberID {
			return tm, nil
		}
	}

	return nil, fmt.Errorf("GetTeamMemberForTeamGroup: %v", errors.New("TeamMember with given ID not found!"))
}

func (db *dbClient) GetAllTeamGroups() ([]*models.TeamGroup, error) {
	teamGroup, err := db.GetTeamGroupsWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllTeamGroups: %v", err)
	}

	return teamGroup, nil
}

func (db *dbClient) GetTeamGroup(teamGroupID string) (*models.TeamGroup, error) {
	selectParams := make(map[string]interface{})
	tg := models.TeamGroup{}
	v := reflect.ValueOf(tg)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return nil, fmt.Errorf("GetTeamGroup: %v", err)
	}
	selectParams[columnName] = teamGroupID

	teamGroups, err := db.GetTeamGroupsWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamGroup: %v", err)
	}

	var teamGroup *models.TeamGroup
	if teamGroups == nil || len(teamGroups) <= 0 {
		return nil, fmt.Errorf("GetTeamGroup: %v", errors.New("Team Group with given ID not found!"))
	} else {
		teamGroup = teamGroups[0]
	}

	return teamGroup, nil
}

func (db *dbClient) GetTeamGroupsWithFilters(searchParams map[string]interface{}) ([]*models.TeamGroup, error) {
	p := models.TeamGroup{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamGroupsWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetTeamGroupsWithFilters: %v", err)
	}

	teamGroups, err := db.GetTeamGroupsFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetTeamGroupsWithFilters: %v", err)
	}

	return teamGroups, nil
}

func (db *dbClient) GetTeamGroupsFromRows(rows *sql.Rows) ([]*models.TeamGroup, error) {
	teamGroups := make([]*models.TeamGroup, 0)
	for rows.Next() {
		tg := models.TeamGroup{}

		var workspaceID sql.NullString

		err := rows.Scan(&tg.ID, &tg.Name, &workspaceID)
		if err != nil {
			return nil, fmt.Errorf("GetTeamGroupsFromRows: %v", err)
		}

		if workspaceID.Valid {
			tg.Workspace, err = db.GetWorkspace(workspaceID.String)
			if err != nil {
				return nil, fmt.Errorf("GetTeamGroupsFromRows: %v", err)
			}
		}

		tg.TeamMembers, err = db.GetTeamGroupTeamMembers(tg.ID)
		if err != nil {
			return nil, fmt.Errorf("GetTeamGroupsFromRows: %v", err)
		}

		teamGroups = append(teamGroups, &tg)
	}

	return teamGroups, nil
}

func (db *dbClient) GetTeamGroupTeamMembers(teamGroupID string) ([]*models.TeamMember, error) {
	searchParams := make(map[string]interface{})
	searchParams["team_group_id"] = teamGroupID

	rows, err := db.GetValuesFromCompositeTable(TEAM_GROUP_TEAM_MEMBER, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamGroupTeamMembers: %v", err)
	}

	teamMembers := make([]*models.TeamMember, 0)
	for rows.Next() {
		teamMemberID := ""

		err := rows.Scan(&teamGroupID, &teamMemberID)
		if err != nil {
			return nil, fmt.Errorf("GetTeamGroupTeamMembers: %v", err)
		}

		tm, err := db.GetTeamMember(teamMemberID)
		if err != nil {
			return nil, fmt.Errorf("GetTeamGroupTeamMembers: %v", err)
		}

		teamMembers = append(teamMembers, tm)
	}

	return teamMembers, nil
}

func (db *dbClient) UpdateTeamGroup(teamGroupID string, updates map[string]interface{}) (*models.TeamGroup, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.TeamGroup{}, teamGroupID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamGroup: %v", err)
	}

	_, err = db.RunUpdateQuery(updateQuery)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamGroup: %v", err)
	}

	teamGroup, err := db.GetTeamGroup(teamGroupID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamGroup: %v", err)
	}

	if v, ok := updates["team_group_team_members"]; ok {
		teamMembers := v.([]*models.TeamMember)
		err = db.UpdateTeamGroupTeamMembers(teamGroupID, teamMembers)
		if err != nil {
			return nil, fmt.Errorf("UpdateProject: %v", err)
		}
	}

	return teamGroup, nil
}

func (db *dbClient) UpdateTeamGroupTeamMembers(teamGroupID string, teamMembers []*models.TeamMember) error {
	deleteParams := make(map[string]interface{})
	_, err := db.DeleteValuesFromCompositeTable(TEAM_GROUP_TEAM_MEMBER, deleteParams)
	if err != nil {
		return fmt.Errorf("UpdateTeamGroupTeamMembers: %v", err)
	}

	err = db.AddTeamGroupTeamMembers(teamGroupID, teamMembers)
	if err != nil {
		return fmt.Errorf("UpdateTeamGroupTeamMembers: %v", err)
	}

	return nil
}

func (db *dbClient) DeleteTeamGroup(teamGroupID string) error {
	deleteParams := make(map[string]interface{})
	c := models.TeamGroup{}
	v := reflect.ValueOf(c)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return fmt.Errorf("DeleteTeamGroup: %v", err)
	}

	deleteParams[columnName] = teamGroupID

	deleteQuery, err := db.GetDeleteQueryForStruct(c, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteTeamGroup: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteTeamGroup: %v", err)
	}

	deleteParamsForColumns := make(map[string]interface{})
	deleteParamsForColumns["team_group_id"] = teamGroupID

	_, err = db.DeleteValuesFromCompositeTable(TEAM_GROUP_TEAM_MEMBER, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteTeamMemebersForTeamGroup: %v", err)
	}

	_, err = db.DeleteValuesFromCompositeTable(PROJECT_TEAM_GROUP, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteProjectsForTeamGroup: %v", err)
	}

	return nil
}
