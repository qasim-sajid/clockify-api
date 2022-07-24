package main

import (
	"fmt"
	"time"

	"github.com/qasim-sajid/clockify-api/dbhandler"
	"github.com/qasim-sajid/clockify-api/models"
)

// Main function
func main() {
	dbC, err := dbhandler.NewDBClient("DB")
	if err != nil {
		panic(err)
	}

	// tasks, err := dbC.GetAllTasks()
	// if err != nil {
	// 	panic(err)
	// }

	// for i := 0; i < len(tasks); i++ {
	// 	fmt.Println(tasks[i])
	// }

	task := models.Task{}
	task.ID = "ID8"
	task.Billable = true
	task.Description = "Description8"
	task.StartTime = time.Now()
	task.EndTime = time.Now()
	task.Date = time.Now()
	task.Project = &models.Project{ID: "Project2"}
	_, status, err := dbC.AddTask(&task)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status: ", status)
}
