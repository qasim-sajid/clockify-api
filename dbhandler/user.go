package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddUser(user *models.User) (*models.User, int, error) {
	if status, err := db.CheckForDuplicateUser(user.Email); err != nil {
		return nil, status, fmt.Errorf("AddUser: %v", err)
	}

	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("Unable to generate _ID")
	}
	user.ID = fmt.Sprintf("u_%v", id)

	insertQuery, err := db.GetInsertQuery(*user)
	if err != nil {
		return nil, -1, fmt.Errorf("AddUser: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddUser: %v", err)
	}

	return user, http.StatusOK, nil
}

func (db *dbClient) GetAllUsers() ([]*models.User, error) {
	users, err := db.GetUsersWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}

	return users, nil
}

func (db *dbClient) GetUser(userID string) (*models.User, error) {
	selectParams := make(map[string]interface{})

	selectParams["_id"] = userID

	users, err := db.GetUsersWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetUser: %v", err)
	}

	var user *models.User
	if users == nil || len(users) <= 0 {
		return nil, nil //fmt.Errorf("GetUser: %v", errors.New("User with given ID not found!"))
	} else {
		user = users[0]
	}

	return user, nil
}

func (db *dbClient) GetUserWithEmail(email string) (*models.User, error) {
	selectParams := make(map[string]interface{})

	selectParams["email"] = email

	users, err := db.GetUsersWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetUser: %v", err)
	}

	var user *models.User
	if users == nil || len(users) <= 0 {
		return nil, nil //fmt.Errorf("GetUser: %v", errors.New("User with given ID not found!"))
	} else {
		user = users[0]
	}

	return user, nil
}

func (db *dbClient) GetUsersWithFilters(searchParams map[string]interface{}) ([]*models.User, error) {
	p := models.User{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetUsersWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetUsersWithFilters: %v", err)
	}

	users, err := db.GetUsersFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetUsersWithFilters: %v", err)
	}

	return users, nil
}

func (db *dbClient) GetUsersFromRows(rows *sql.Rows) ([]*models.User, error) {
	users := make([]*models.User, 0)
	for rows.Next() {
		u := models.User{}

		err := rows.Scan(&u.ID, &u.Email, &u.Name)

		if err != nil {
			return nil, fmt.Errorf("GetUsersFromRows: %v", err)
		}

		users = append(users, &u)
	}

	return users, nil
}

func (db *dbClient) UpdateUser(userID string, updates map[string]interface{}) (*models.User, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.User{}, userID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateUser: %v", err)
	}

	if len(updates) > 0 {
		_, err = db.RunUpdateQuery(updateQuery)
		if err != nil {
			return nil, fmt.Errorf("UpdateUser: %v", err)
		}
	}

	user, err := db.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("UpdateUser: %v", err)
	}

	return user, nil
}

func (db *dbClient) DeleteUser(userID string) error {
	deleteParams := make(map[string]interface{})

	deleteParams["_id"] = userID

	deleteQuery, err := db.GetDeleteQueryForStruct(models.User{}, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteUser: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteUser: %v", err)
	}

	return nil
}

func (db *dbClient) CheckForDuplicateUser(email string) (int, error) {
	searchParams := make(map[string]interface{})
	searchParams["email"] = email
	user, err := db.GetUsersWithFilters(searchParams)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("CheckForDuplicateUser: %v", err)
	}
	if user != nil {
		return http.StatusBadRequest, fmt.Errorf("CheckForDuplicateUser: %v", errors.New("User with this email already exists!"))
	}

	return http.StatusOK, nil
}
