package models

import (
	"time"
)

// Task defines task object
type Task struct {
	ID          string    `json:"_id"`
	Description string    `json:"description"`
	Billable    bool      `json:"billable"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Date        time.Time `json:"date"`
	IsActive    bool      `json:"is_active"`

	Project string   `json:"project_id"`
	Tags    []string `json:"tags"`
}
