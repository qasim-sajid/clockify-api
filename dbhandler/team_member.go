package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddTeamMember(teamMember *models.TeamMember) (*models.TeamMember, int, error) {
	insertQuery, err := db.GetInsertQueryForStruct(teamMember)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamMember: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamMember: %v", err)
	}

	return teamMember, http.StatusOK, nil
}

func (db *dbClient) GetAllTeamMembers() ([]*models.TeamMember, error) {
	teamMembers, err := db.GetTeamMembersWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllTeamMembers: %v", err)
	}

	return teamMembers, nil
}

func (db *dbClient) GetTeamMember(teamMemberID string) (*models.TeamMember, error) {
	selectParams := make(map[string]interface{})
	c := models.TeamMember{}
	v := reflect.ValueOf(c)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return nil, fmt.Errorf("GetTeamMember: %v", err)
	}
	selectParams[columnName] = teamMemberID

	teamMembers, err := db.GetTeamMembersWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamMember: %v", err)
	}

	var teamMember *models.TeamMember
	if teamMembers == nil {
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

		err := rows.Scan(&tm.ID, &tm.BillableRate)

		if err != nil {
			return nil, fmt.Errorf("GetTeamMembersFromRows: %v", err)
		}

		teamMembers = append(teamMembers, &tm)
	}

	return teamMembers, nil
}

func (db *dbClient) UpdateTeamMember(teamMemberID string, updates map[string]interface{}) (*models.TeamMember, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.TeamMember{}, teamMemberID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamMember: %v", err)
	}

	_, err = db.RunUpdateQuery(updateQuery)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamMember: %v", err)
	}

	teamMember, err := db.GetTeamMember(teamMemberID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamMember: %v", err)
	}

	return teamMember, nil
}

func (db *dbClient) DeleteTeamMember(teamMemberID string) error {
	deleteParams := make(map[string]interface{})
	c := models.TeamMember{}
	v := reflect.ValueOf(c)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return fmt.Errorf("DeleteTeamMember: %v", err)
	}

	deleteParams[columnName] = teamMemberID

	deleteQuery, err := db.GetDeleteQueryForStruct(c, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteTeamMember: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteTeamMember: %v", err)
	}

	return nil
}
