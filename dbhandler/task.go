package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/models"
)

func (db *dbClient) AddTask(task *models.Task) (*models.Task, int, error) {
	id := uuid.New().String()
	if id == "" {
		return nil, http.StatusInternalServerError, errors.New("Unable to generate _ID")
	}
	task.ID = fmt.Sprintf("t_%v", id)

	insertQuery, err := db.GetInsertQuery(*task)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTask: %v", err)
	}

	_, err = db.RunInsertQuery(insertQuery)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTask: %v", err)
	}

	err = db.AddTaskTags(task.ID, task.Tags)
	if err != nil {
		return nil, -1, fmt.Errorf("AddTask: %v", err)
	}

	return task, http.StatusOK, nil
}

func (db *dbClient) AddTaskTags(taskID string, tags []string) error {
	if tags == nil {
		return nil
	}

	for _, t := range tags {
		valuesMap := make(map[string]interface{})
		valuesMap["task_id"] = taskID
		valuesMap["tag_id"] = t

		//Check if value already exists
		_, err := db.GetTagForTask(taskID, t)
		if err != nil {
			//If value doesn't exist then insert it
			insertQuery, err := db.GetInsertQueryForCompositeTable(TASK_TAG, valuesMap)
			if err != nil {
				return fmt.Errorf("AddTaskTags: %v", err)
			}

			_, err = db.RunInsertQuery(insertQuery)
			if err != nil {
				return fmt.Errorf("AddTaskTags: %v", err)
			}
		}
	}

	return nil
}

func (db *dbClient) GetTagForTask(taskID, tagID string) (string, error) {
	tags, err := db.GetTaskTags(taskID)
	if err != nil {
		return "", fmt.Errorf("GetTagForTask: %v", err)
	}

	for _, t := range tags {
		if t == tagID {
			return t, nil
		}
	}

	return "", fmt.Errorf("GetTagForTask: %v", errors.New("Tag with given ID not found!"))
}

func (db *dbClient) GetAllTasks() ([]*models.Task, error) {
	tasks, err := db.GetTasksWithFilters(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("GetAllTasks: %v", err)
	}

	return tasks, nil
}

func (db *dbClient) GetTask(taskID string) (*models.Task, error) {
	selectParams := make(map[string]interface{})

	selectParams["_id"] = taskID

	tasks, err := db.GetTasksWithFilters(selectParams)
	if err != nil {
		return nil, fmt.Errorf("GetTask: %v", err)
	}

	var task *models.Task
	if tasks == nil || len(tasks) <= 0 {
		return nil, fmt.Errorf("GetTask: %v", errors.New("Task with given ID not found!"))
	} else {
		task = tasks[0]
	}

	return task, nil
}

func (db *dbClient) GetTasksWithFilters(searchParams map[string]interface{}) ([]*models.Task, error) {
	p := models.Task{}

	selectQuery, err := db.GetSelectQueryForStruct(p, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTasksWithFilters: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetTasksWithFilters: %v", err)
	}

	tasks, err := db.GetTasksFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetTasksWithFilters: %v", err)
	}

	return tasks, nil
}

func (db *dbClient) GetTasksFromRows(rows *sql.Rows) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)
	for rows.Next() {
		t := models.Task{}

		var projectID sql.NullString
		startTime := ""
		endTime := ""
		date := ""

		err := rows.Scan(&t.ID, &t.Description, &t.Billable, &startTime, &endTime, &date, &t.IsActive, &projectID)

		if err != nil {
			return nil, fmt.Errorf("GetTasksFromRows: %v", err)
		}

		t.StartTime, err = time.Parse(conf.TIME_LAYOUT, startTime)
		if err != nil {
			return nil, fmt.Errorf("GetTasksFromRows: %v", err)
		}

		t.EndTime, err = time.Parse(conf.TIME_LAYOUT, endTime)
		if err != nil {
			return nil, fmt.Errorf("GetTasksFromRows: %v", err)
		}

		t.Date, err = time.Parse(conf.TIME_LAYOUT, date)
		if err != nil {
			return nil, fmt.Errorf("GetTasksFromRows: %v", err)
		}

		t.Project = projectID.String

		t.Tags, err = db.GetTaskTags(t.ID)
		if err != nil {
			return nil, fmt.Errorf("GetTasksFromRows: %v", err)
		}

		tasks = append(tasks, &t)
	}

	return tasks, nil
}

func (db *dbClient) GetTaskTags(taskID string) ([]string, error) {
	searchParams := make(map[string]interface{})
	searchParams["task_id"] = taskID

	selectQuery, err := db.GetSelectQueryForCompositeTable(TASK_TAG, searchParams)
	if err != nil {
		return nil, fmt.Errorf("GetTaskTags: %v", err)
	}

	rows, err := db.RunSelectQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("GetTaskTags: %v", err)
	}

	tags := make([]string, 0)
	for rows.Next() {
		tagID := ""

		err := rows.Scan(&taskID, &tagID)
		if err != nil {
			return nil, fmt.Errorf("GetTaskTags: %v", err)
		}

		tags = append(tags, tagID)
	}

	return tags, nil
}

func (db *dbClient) UpdateTask(taskID string, updates map[string]interface{}) (*models.Task, error) {
	if v, ok := updates["tags"]; ok {
		tagIDs := strings.Split(v.(string), ",")

		err := db.UpdateTaskTags(taskID, tagIDs)
		if err != nil {
			return nil, fmt.Errorf("UpdateTask: %v", err)
		}
		delete(updates, "tags")
	}

	updateQuery, err := db.GetUpdateQueryForStruct(models.Task{}, taskID, updates)
	if err != nil {
		return nil, fmt.Errorf("UpdateTask: %v", err)
	}

	if len(updates) > 0 {
		_, err = db.RunUpdateQuery(updateQuery)
		if err != nil {
			return nil, fmt.Errorf("UpdateTask: %v", err)
		}
	}

	task, err := db.GetTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTask: %v", err)
	}

	return task, nil
}

func (db *dbClient) UpdateTaskTags(taskID string, tags []string) error {
	deleteParams := make(map[string]interface{})
	deleteParams["task_id"] = taskID
	_, err := db.DeleteValuesFromCompositeTable(TASK_TAG, deleteParams)
	if err != nil {
		return fmt.Errorf("UpdateTaskTags: %v", err)
	}

	err = db.AddTaskTags(taskID, tags)
	if err != nil {
		return fmt.Errorf("UpdateTaskTags: %v", err)
	}

	return nil
}

func (db *dbClient) DeleteTask(taskID string) error {
	deleteParams := make(map[string]interface{})

	deleteParams["_id"] = taskID

	deleteQuery, err := db.GetDeleteQueryForStruct(models.Task{}, deleteParams)
	if err != nil {
		return fmt.Errorf("DeleteTask: %v", err)
	}

	_, err = db.RunDeleteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("DeleteTask: %v", err)
	}

	deleteParamsForColumns := make(map[string]interface{})
	deleteParamsForColumns["task_id"] = taskID

	_, err = db.DeleteValuesFromCompositeTable(TASK_TAG, deleteParamsForColumns)
	if err != nil {
		return fmt.Errorf("DeleteTagsForTask: %v", err)
	}

	return nil
}
