package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddTag(tag *models.Tag) (*models.Tag, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("unable to generate id")
	}
	tag.ID = fmt.Sprintf("t_%v", id)

	insertQuery, err := db.GetInsertQuery(*tag)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTag: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTag: %v", err)
	}

	return tag, http.StatusOK, nil
}

func (db *dbClient) GetAllTags() ([]*models.Tag, error) {
	tags, err := db.GetTagsWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllTags: %v", err)
	}

	return tags, nil
}

func (db *dbClient) GetTag(tagID string) (*models.Tag, error) {
	selectParams := make(map[string]interface{})

	selectParams["_id"] = tagID

	tags, err := db.GetTagsWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetTag: %v", err)
	}

	var tag *models.Tag
	if tags == nil || len(tags) <= 0 {
		return nil, fmt.Errorf("GetTag: %v", errors.New("tag with given id not found"))
	} else {
		tag = tags[0]
	}

	return tag, nil
}

func (db *dbClient) GetTagsWithFilters(searchParams map[string]interface{}) ([]*models.Tag, error) {
	p := models.Tag{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTagsWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetTagsWithFilters: %v", err)
	}

	tags, err := db.GetTagsFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetTagsWithFilters: %v", err)
	}

	return tags, nil
}

func (db *dbClient) GetTagsFromRows(rows *sql.Rows) ([]*models.Tag, error) {
	tags := make([]*models.Tag, 0)
	for rows.Next() {
		t := models.Tag{}

		err := rows.Scan(&t.ID, &t.Name)

		if err != nil {
			return nil, fmt.Errorf("GetTagsFromRows: %v", err)
		}

		tags = append(tags, &t)
	}

	return tags, nil
}

func (db *dbClient) UpdateTag(tagID string, updates map[string]interface{}) (*models.Tag, error) {
	updateQuery, err := db.GetUpdateQueryForStruct(models.Tag{}, tagID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateTag: %v", err)
	}

	if len(updates) > 0 {
		_, err = db.RunUpdateQuery(updateQuery)
		if err != nil {
			return nil, fmt.Errorf("UpdateTag: %v", err)
		}
	}

	tag, err := db.GetTag(tagID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTag: %v", err)
	}

	return tag, nil
}

func (db *dbClient) DeleteTag(tagID string) error {
	deleteParams := make(map[string]interface{})

	deleteParams["_id"] = tagID

	deleteQuery, err := db.GetDeleteQueryForStruct(models.Tag{}, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteTag: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteTag: %v", err)
	}

	deleteParamsForColumns := make(map[string]interface{})
	deleteParamsForColumns["tag_id"] = tagID

	_, err = db.DeleteValuesFromCompositeTable(TASK_TAG, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteTasksForTag: %v", err)
	}

	return nil
}
