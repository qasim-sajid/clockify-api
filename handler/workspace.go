package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddWorkspace(c *gin.Context, h *Handler, origin *models.User) {
	workspace := &models.Workspace{}

	workspace.Name = c.Query("name")

	workspace, _, err := h.DB.AddWorkspace(workspace)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Workspace with _id = %s added!", workspace.ID)})
	}
}

func GetAllWorkspaces(c *gin.Context, h *Handler, origin *models.User) {
	workspaces, err := h.DB.GetAllWorkspaces()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

func GetWorkspace(c *gin.Context, h *Handler, origin *models.User) {
	workspaceID := c.Param("workspace_id")
	workspace, err := h.DB.GetWorkspace(workspaceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, workspace)
	}
}

func UpdateWorkspace(c *gin.Context, h *Handler, origin *models.User) {
	updates := make(map[string]interface{})
	workspaceID := c.Param("workspace_id")
	for k, v := range c.Request.URL.Query() {
		if len(v) > 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
			return
		} else {
			updates[k] = v[0]
		}
	}

	_, err := h.DB.UpdateWorkspace(workspaceID, updates)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Workspace with _id = %s updated!", workspaceID)})
	}
}

func DeleteWorkspace(c *gin.Context, h *Handler, origin *models.User) {
	workspaceID := c.Param("workspace_id")
	err := h.DB.DeleteWorkspace(workspaceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Workspace with _id = %s deleted!", workspaceID)})
	}
}
