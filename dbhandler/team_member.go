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

func (db *dbClient) AddTeamMember(teamMember *models.TeamMember) (*models.TeamMember, int, error) {
	err := db.checkForDuplicateTeamMember(teamMember)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("AddTeamMember: %v", err)
	}

	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("AddTeamMember: %v", errors.New("Unable to generate _ID"))
	}
	teamMember.ID = fmt.Sprintf("tm_%v", id)

	insertQuery, err := db.GetInsertQuery(*teamMember)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamMember: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamMember: %v", err)
	}

	return teamMember, http.StatusOK, nil
}

func (db *dbClient) checkForDuplicateTeamMember(teamMember *models.TeamMember) error {
	searchParams := make(map[string]interface{})
	searchParams["user_email"] = teamMember.User
	teamMembers, _ := db.GetTeamMembersWithFilters(searchParams)
	if teamMembers != nil && len(teamMembers) > 0 {
		return errors.New("TeamMember with this user email already exists!")
	}

	return nil
}

func (db *dbClient) GetAllTeamMembers() ([]*models.TeamMember, error) {
	teamMembers, err := db.GetTeamMembersWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllTeamMembers: %v", err)
	}

	return teamMembers, nil
}

func (db *dbClient) AddTeamMemberTeamGroups(teamMemberID string, teamGroups []string) error {
	if teamGroups == nil {
		return nil
	}

	for _, tg := range teamGroups {
		valuesMap := make(map[string]interface{})
		valuesMap["team_group_id"] = tg
		valuesMap["team_member_id"] = teamMemberID

		//Check if value already exists
		_, err := db.GetTeamGroupForTeamMember(teamMemberID, tg)
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

func (db *dbClient) GetTeamGroupForTeamMember(teamMemberID, teamGroupID string) (string, error) {
	teamGroups, err := db.GetTeamMemberTeamGroups(teamMemberID)
	if err != nil {
		return "", fmt.Errorf("GetTeamGroupForTeamMember: %v", err)
	}

	for _, tg := range teamGroups {
		if tg == teamGroupID {
			return tg, nil
		}
	}

	return "", fmt.Errorf("GetTeamGroupForTeamMember: %v", errors.New("TeamGroup with given ID not found!"))
}

func (db *dbClient) GetTeamMember(teamMemberID string) (*models.TeamMember, error) {
	selectParams := make(map[string]interface{})

	selectParams["_id"] = teamMemberID

	teamMembers, err := db.GetTeamMembersWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMember: %v", err)
	}

	var teamMember *models.TeamMember
	if teamMembers == nil || len(teamMembers) <= 0 {
		return nil, fmt.Errorf("GetTeamMember: %v", errors.New("TeamMember with given ID not found!"))
	} else {
		teamMember = teamMembers[0]
	}

	return teamMember, nil
}

func (db *dbClient) GetTeamMembersWithFilters(searchParams map[string]interface{}) ([]*models.TeamMember, error) {
	p := models.TeamMember{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMembersWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMembersWithFilters: %v", err)
	}

	teamMembers, err := db.GetTeamMembersFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMembersWithFilters: %v", err)
	}

	return teamMembers, nil
}

func (db *dbClient) GetTeamMembersFromRows(rows *sql.Rows) ([]*models.TeamMember, error) {
	teamMembers := make([]*models.TeamMember, 0)
	for rows.Next() {
		tm := models.TeamMember{}

		var workspaceID sql.NullString
		var userEmail sql.NullString
		var teamRoleID sql.NullString

		err := rows.Scan(&tm.ID, &tm.BillableRate, &workspaceID, &userEmail, &teamRoleID)
		if err != nil {
			return nil, fmt.Errorf("GetTeamMembersFromRows: %v", err)
		}

		if workspaceID.Valid {
			tm.Workspace = workspaceID.String
		}

		if userEmail.Valid {
			tm.User = userEmail.String
		}

		if teamRoleID.Valid {
			tm.TeamRole = teamRoleID.String
		}

		tm.TeamGroups, err = db.GetTeamMemberTeamGroups(tm.ID)
		if err != nil {
			return nil, fmt.Errorf("GetTeamMemberTeamGroups: %v", err)
		}

		teamMembers = append(teamMembers, &tm)
	}

	return teamMembers, nil
}

func (db *dbClient) GetTeamMemberTeamGroups(teamMemberID string) ([]string, error) {
	searchParams := make(map[string]interface{})
	searchParams["team_member_id"] = teamMemberID

	rows, err := db.GetValuesFromCompositeTable(TEAM_GROUP_TEAM_MEMBER, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMemberTeamGroups: %v", err)
	}

	teamGroups := make([]string, 0)
	for rows.Next() {
		teamGroupID := ""

		err := rows.Scan(&teamGroupID, &teamMemberID)
		if err != nil {
			return nil, fmt.Errorf("GetTeamMemberTeamGroups: %v", err)
		}

		teamGroups = append(teamGroups, teamGroupID)
	}

	return teamGroups, nil
}

func (db *dbClient) UpdateTeamMember(teamMemberID string, updates map[string]interface{}) (*models.TeamMember, error) {
	if v, ok := updates["team_groups"]; ok {
		teamGroups := strings.Split(v.(string), ",")

		err := db.UpdateTeamMemberTeamGroups(teamMemberID, teamGroups)
		if err != nil {
			return nil, fmt.Errorf("UpdateTeamMember: %v", err)
		}
		delete(updates, "team_groups")
	}

	updateQuery, err := db.GetUpdateQueryForStruct(models.TeamMember{}, teamMemberID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamMember: %v", err)
	}

	if len(updates) > 0 {
		_, err = db.RunUpdateQuery(updateQuery)
		if err != nil {
			return nil, fmt.Errorf("UpdateTeamMember: %v", err)
		}
	}

	teamMember, err := db.GetTeamMember(teamMemberID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamMember: %v", err)
	}

	return teamMember, nil
}

func (db *dbClient) UpdateTeamMemberTeamGroups(teamMemberID string, teamGroups []string) error {
	deleteParams := make(map[string]interface{})
	deleteParams["team_member_id"] = teamMemberID
	_, err := db.DeleteValuesFromCompositeTable(TEAM_GROUP_TEAM_MEMBER, deleteParams)
	if err != nil {
		return fmt.Errorf("UpdateTeamMemberTeamGroups: %v", err)
	}

	err = db.AddTeamMemberTeamGroups(teamMemberID, teamGroups)
	if err != nil {
		return fmt.Errorf("UpdateTeamMemberTeamGroups: %v", err)
	}

	return nil
}

func (db *dbClient) DeleteTeamMember(teamMemberID string) error {
	deleteParams := make(map[string]interface{})

	deleteParams["_id"] = teamMemberID

	deleteQuery, err := db.GetDeleteQueryForStruct(models.TeamMember{}, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteTeamMember: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteTeamMember: %v", err)
	}

	deleteParamsForColumns := make(map[string]interface{})
	deleteParamsForColumns["team_member_id"] = teamMemberID

	_, err = db.DeleteValuesFromCompositeTable(PROJECT_TEAM_MEMBER, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteProjectsForTeamMember: %v", err)
	}

	_, err = db.DeleteValuesFromCompositeTable(TEAM_GROUP_TEAM_MEMBER, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteTeamGroupsForTeamMember: %v", err)
	}

	return nil
}
