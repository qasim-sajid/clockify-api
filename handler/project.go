package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddProject(c *gin.Context, h *Handler, origin *models.User) {
	project := &models.Project{}

	project.Name = c.Query("name")
	project.ColorTag = c.Query("color_tag")
	var err error

	project.IsPublic, err = strconv.ParseBool(c.Query("is_public"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	project.TrackedHours, err = strconv.ParseFloat(c.Query("tracked_hours"), 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	project.TrackedAmount, err = strconv.ParseFloat(c.Query("tracked_amount"), 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	project.ProgressPercentage, err = strconv.ParseFloat(c.Query("progress_percentage"), 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	project.Client = c.Query("client_id")

	project.Workspace = c.Query("workspace_id")

	project, _, err = h.DB.AddProject(project)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Project with _id = %s added!", project.ID)})
	}
}

func GetAllProjects(c *gin.Context, h *Handler, origin *models.User) {
	projects, err := h.DB.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, projects)
}

func GetProject(c *gin.Context, h *Handler, origin *models.User) {
	projectID := c.Param("project_id")
	project, err := h.DB.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, project)
	}
}

func UpdateProject(c *gin.Context, h *Handler, origin *models.User) {
	updates := make(map[string]interface{})
	projectID := c.Param("project_id")
	for k, v := range c.Request.URL.Query() {
		if len(v) > 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
		} else {
			updates[k] = v[0]
		}
	}

	_, err := h.DB.UpdateProject(projectID, updates)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Project with _id = %s updated!", projectID)})
	}
}

func DeleteProject(c *gin.Context, h *Handler, origin *models.User) {
	projectID := c.Param("project_id")
	err := h.DB.DeleteProject(projectID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Project with _id = %s deleted!", projectID)})
	}
}
