package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddTeamMember(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamMember := &models.TeamMember{}

		var err error

		teamMember.BillableRate, err = strconv.ParseFloat(c.Query("billable_rate"), 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		teamMember.Workspace = c.Query("workspace_id")

		teamMember.User = c.Query("user_email")

		teamMember.TeamRole = c.Query("team_role_id")

		teamMember, _, err = h.DB.AddTeamMember(teamMember)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamMember with _id = %s added!", teamMember.ID)})
		}
	}
}

func GetAllTeamMembers(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamMembers, err := h.DB.GetAllTeamMembers()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, teamMembers)
	}
}

func GetTeamMember(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamMemberID := c.Param("team_member_id")
		teamMember, err := h.DB.GetTeamMember(teamMemberID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, teamMember)
		}
	}
}

func UpdateTeamMember(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		updates := make(map[string]interface{})
		teamMemberID := c.Param("team_member_id")
		for k, v := range c.Request.URL.Query() {
			if len(v) > 1 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
			} else {
				updates[k] = v[0]
			}
		}

		_, err := h.DB.UpdateTeamMember(teamMemberID, updates)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamMember with _id = %s updated!", teamMemberID)})
		}
	}
}

func DeleteTeamMember(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamMemberID := c.Param("team_member_id")
		err := h.DB.DeleteTeamMember(teamMemberID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("TeamMember with _id = %s deleted!", teamMemberID)})
		}
	}
}
