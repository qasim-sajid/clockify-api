package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddTeamGroup(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamGroup := &models.TeamGroup{}

		teamGroup.Name = c.Query("name")

		var err error

		teamGroup.Workspace = c.Query("workspace_id")

		teamGroup, _, err = h.DB.AddTeamGroup(teamGroup)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamGroup with _id = %s added!", teamGroup.ID)})
		}
	}
}

func GetAllTeamGroups(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamGroups, err := h.DB.GetAllTeamGroups()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, teamGroups)
	}
}

func GetTeamGroup(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamGroupID := c.Param("team_group_id")
		teamGroup, err := h.DB.GetTeamGroup(teamGroupID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, teamGroup)
		}
	}
}

func UpdateTeamGroup(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		updates := make(map[string]interface{})
		teamGroupID := c.Param("team_group_id")
		for k, v := range c.Request.URL.Query() {
			if len(v) > 1 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
			} else {
				updates[k] = v[0]
			}
		}

		_, err := h.DB.UpdateTeamGroup(teamGroupID, updates)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamGroup with _id = %s updated!", teamGroupID)})
		}
	}
}

func DeleteTeamGroup(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamGroupID := c.Param("team_group_id")
		err := h.DB.DeleteTeamGroup(teamGroupID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamGroup with _id = %s deleted!", teamGroupID)})
		}
	}
}
