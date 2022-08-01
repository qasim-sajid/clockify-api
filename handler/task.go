package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddTask(c *gin.Context, h *Handler, origin *models.User) {
	task := &models.Task{}

	task.Description = c.Query("description")
	var err error

	task.Billable, err = strconv.ParseBool(c.Query("billable"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	task.StartTime, err = time.Parse(conf.TIME_LAYOUT, c.Query("start_time"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	task.EndTime, err = time.Parse(conf.TIME_LAYOUT, c.Query("end_time"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	task.Date, err = time.Parse(conf.TIME_LAYOUT, c.Query("date"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	task.IsActive, err = strconv.ParseBool(c.Query("is_active"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	task.Project = c.Query("project_id")

	tags := strings.Split(c.Query("tags"), ",")
	if len(tags) > 0 && tags[0] != "" {
		task.Tags = tags
	}

	task, _, err = h.DB.AddTask(task)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Task with _id = %s added!", task.ID)})
	}
}

func GetAllTasks(c *gin.Context, h *Handler, origin *models.User) {
	tasks, err := h.DB.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context, h *Handler, origin *models.User) {
	taskID := c.Param("task_id")
	task, err := h.DB.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, task)
	}
}

func UpdateTask(c *gin.Context, h *Handler, origin *models.User) {
	updates := make(map[string]interface{})
	taskID := c.Param("task_id")
	for k, v := range c.Request.URL.Query() {
		if len(v) > 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
		} else {
			updates[k] = v[0]
		}
	}

	_, err := h.DB.UpdateTask(taskID, updates)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Task with _id = %s updated!", taskID)})
	}
}

func DeleteTask(c *gin.Context, h *Handler, origin *models.User) {
	taskID := c.Param("task_id")
	err := h.DB.DeleteTask(taskID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Task with _id = %s deleted!", taskID)})
	}
}
