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

func (db *dbClient) AddWorkspace(workspace *models.Workspace) (*models.Workspace, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("Unable to generate _ID")
	}
	workspace.ID = fmt.Sprintf("w_%v", id)

	insertQuery, err := db.GetInsertQuery(*workspace)
	if err != nil {
		return nil, -1, fmt.Errorf("AddWorkspace: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddWorkspace: %v", err)
	}

	return workspace, http.StatusOK, nil
}

func (db *dbClient) GetAllWorkspaces() ([]*models.Workspace, error) {
	workspaces, err := db.GetWorkspacesWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllWorkspaces: %v", err)
	}

	return workspaces, nil
}

func (db *dbClient) GetWorkspace(workspaceID string) (*models.Workspace, error) {
	selectParams := make(map[string]interface{})
	w := models.Workspace{}
	v := reflect.ValueOf(w)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return nil, fmt.Errorf("GetWorkspace: %v", err)
	}
	selectParams[columnName] = workspaceID

	workspaces, err := db.GetWorkspacesWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetWorkspace: %v", err)
	}

	var workspace *models.Workspace
	if workspaces == nil || len(workspaces) <= 0 {
		return nil, fmt.Errorf("GetWorkspace: %v", errors.New("Workspace with given ID not found!"))
	} else {
		workspace = workspaces[0]
	}

	return workspace, nil
}

func (db *dbClient) GetWorkspacesWithFilters(searchParams map[string]interface{}) ([]*models.Workspace, error) {
	p := models.Workspace{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetWorkspacesWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetWorkspacesWithFilters: %v", err)
	}

	workspaces, err := db.GetWorkspacesFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetWorkspacesWithFilters: %v", err)
	}

	return workspaces, nil
}

func (db *dbClient) GetWorkspacesFromRows(rows *sql.Rows) ([]*models.Workspace, error) {
	workspaces := make([]*models.Workspace, 0)
	for rows.Next() {
		w := models.Workspace{}

		err := rows.Scan(&w.ID, &w.Name)

		if err != nil {
			return nil, fmt.Errorf("GetWorkspacesFromRows: %v", err)
		}

		workspaces = append(workspaces, &w)
	}

	return workspaces, nil
}

func (db *dbClient) UpdateWorkspace(workspaceID string, updates map[string]interface{}) (*models.Workspace, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.Workspace{}, workspaceID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateWorkspace: %v", err)
	}

	_, err = db.RunUpdateQuery(updateQuery)
	if err != nil {
		return nil, fmt.Errorf("UpdateWorkspace: %v", err)
	}

	workspace, err := db.GetWorkspace(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("UpdateWorkspace: %v", err)
	}

	return workspace, nil
}

func (db *dbClient) DeleteWorkspace(workspaceID string) error {
	deleteParams := make(map[string]interface{})
	c := models.Workspace{}
	v := reflect.ValueOf(c)

	columnName, err := db.GetColumnNameForStructField(v.Type().Field(0))
	if err != nil {
		return fmt.Errorf("DeleteWorkspace: %v", err)
	}

	deleteParams[columnName] = workspaceID

	deleteQuery, err := db.GetDeleteQueryForStruct(c, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteWorkspace: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteWorkspace: %v", err)
	}

	return nil
}
