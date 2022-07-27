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

func (db *dbClient) AddTeamRole(teamRole *models.TeamRole) (*models.TeamRole, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("Unable to generate _ID")
	}
	teamRole.ID = fmt.Sprintf("tr_%v", id)

	insertQuery, err := db.GetInsertQuery(*teamRole)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamRole: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTeamRole: %v", err)
	}

	return teamRole, http.StatusOK, nil
}

func (db *dbClient) GetAllTeamRoles() ([]*models.TeamRole, error) {
	teamRoles, err := db.GetTeamRolesWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllTeamRoles: %v", err)
	}

	return teamRoles, nil
}

func (db *dbClient) GetTeamRole(teamRoleID string) (*models.TeamRole, error) {
	selectParams := make(map[string]interface{})
	tr := models.TeamRole{}
	v := reflect.ValueOf(tr)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return nil, fmt.Errorf("GetTeamRole: %v", err)
	}
	selectParams[columnName] = teamRoleID

	teamRoles, err := db.GetTeamRolesWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamRole: %v", err)
	}

	var teamRole *models.TeamRole
	if teamRoles == nil || len(teamRoles) <= 0 {
		return nil, fmt.Errorf("GetTeamRole: %v", errors.New("TeamRole with given ID not found!"))
	} else {
		teamRole = teamRoles[0]
	}

	return teamRole, nil
}

func (db *dbClient) GetTeamRolesWithFilters(searchParams map[string]interface{}) ([]*models.TeamRole, error) {
	p := models.TeamRole{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTeamRolesWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetTeamRolesWithFilters: %v", err)
	}

	teamRoles, err := db.GetTeamRolesFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetTeamRolesWithFilters: %v", err)
	}

	return teamRoles, nil
}

func (db *dbClient) GetTeamRolesFromRows(rows *sql.Rows) ([]*models.TeamRole, error) {
	teamRoles := make([]*models.TeamRole, 0)
	for rows.Next() {
		tr := models.TeamRole{}

		err := rows.Scan(&tr.ID, &tr.Role)

		if err != nil {
			return nil, fmt.Errorf("GetTeamRolesFromRows: %v", err)
		}

		teamRoles = append(teamRoles, &tr)
	}

	return teamRoles, nil
}

func (db *dbClient) UpdateTeamRole(teamRoleID string, updates map[string]interface{}) (*models.TeamRole, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.TeamRole{}, teamRoleID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamRole: %v", err)
	}

	_, err = db.RunUpdateQuery(updateQuery)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamRole: %v", err)
	}

	teamRole, err := db.GetTeamRole(teamRoleID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTeamRole: %v", err)
	}

	return teamRole, nil
}

func (db *dbClient) DeleteTeamRole(teamRoleID string) error {
	deleteParams := make(map[string]interface{})
	c := models.TeamRole{}
	v := reflect.ValueOf(c)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return fmt.Errorf("DeleteTeamRole: %v", err)
	}

	deleteParams[columnName] = teamRoleID

	deleteQuery, err := db.GetDeleteQueryForStruct(c, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteTeamRole: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteTeamRole: %v", err)
	}

	return nil
}
