package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddClient(client *models.Client) (*models.Client, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("Unable to generate _ID")
	}
	client.ID = fmt.Sprintf("c_%v", id)

	insertQuery, err := db.GetInsertQuery(*client)
	if err != nil {
		return nil, -1, fmt.Errorf("AddClient: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddClient: %v", err)
	}

	return client, http.StatusOK, nil
}

func (db *dbClient) GetAllClients() ([]*models.Client, error) {
	clients, err := db.GetClientsWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllProjects: %v", err)
	}

	return clients, nil
}

func (db *dbClient) GetClient(clientID string) (*models.Client, error) {
	selectParams := make(map[string]interface{})

	selectParams["_id"] = clientID

	clients, err := db.GetClientsWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetClient: %v", err)
	}

	var client *models.Client
	if clients == nil || len(clients) <= 0 {
		return nil, fmt.Errorf("GetClient: %v", errors.New("Client with given ID not found!"))
	} else {
		client = clients[0]
	}

	return client, nil
}

func (db *dbClient) GetClientsWithFilters(searchParams map[string]interface{}) ([]*models.Client, error) {
	p := models.Client{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetClientsWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetClientsWithFilters: %v", err)
	}

	clients, err := db.GetClientsFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetClientsWithFilters: %v", err)
	}

	return clients, nil
}

func (db *dbClient) GetClientsFromRows(rows *sql.Rows) ([]*models.Client, error) {
	clients := make([]*models.Client, 0)
	for rows.Next() {
		c := models.Client{}

		err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Note, &c.IsArchived)

		if err != nil {
			return nil, fmt.Errorf("GetClientsFromRows: %v", err)
		}

		clients = append(clients, &c)
	}

	return clients, nil
}

func (db *dbClient) UpdateClient(clientID string, updates map[string]interface{}) (*models.Client, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.Client{}, clientID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateClient: %v", err)
	}

	if len(updates) > 0 {
		_, err = db.RunUpdateQuery(updateQuery)
		if err != nil {
			return nil, fmt.Errorf("UpdateClient: %v", err)
		}
	}

	client, err := db.GetClient(clientID)
	if err != nil {
		return nil, fmt.Errorf("UpdateClient: %v", err)
	}

	return client, nil
}

func (db *dbClient) DeleteClient(clientID string) error {
	deleteParams := make(map[string]interface{})

	deleteParams["_id"] = clientID

	deleteQuery, err := db.GetDeleteQueryForStruct(models.Client{}, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteClient: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteClient: %v", err)
	}

	return nil
}
