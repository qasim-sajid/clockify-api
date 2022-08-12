package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddTeamRole(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamRole := &models.TeamRole{}
		teamRole.Role = c.Query("role")

		teamRole, _, err := h.DB.AddTeamRole(teamRole)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamRole with _id = %s added!", teamRole.ID)})
		}
	}
}

func GetAllTeamRoles(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamRoles, err := h.DB.GetAllTeamRoles()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, teamRoles)
	}
}

func GetTeamRole(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamRoleID := c.Param("team_role_id")
		teamRole, err := h.DB.GetTeamRole(teamRoleID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, teamRole)
		}
	}
}

func UpdateTeamRole(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		updates := make(map[string]interface{})
		teamRoleID := c.Param("team_role_id")
		for k, v := range c.Request.URL.Query() {
			if len(v) > 1 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
			} else {
				updates[k] = v[0]
			}
		}

		_, err := h.DB.UpdateTeamRole(teamRoleID, updates)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamRole with _id = %s updated!", teamRoleID)})
		}
	}
}

func DeleteTeamRole(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamRoleID := c.Param("team_role_id")
		err := h.DB.DeleteTeamRole(teamRoleID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamRole with _id = %s deleted!", teamRoleID)})
		}
	}
}
