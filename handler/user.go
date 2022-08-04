package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddUser(c *gin.Context, h *Handler, origin *models.User) {
	user := &models.User{}
	user.Name = c.Query("name")
	user.Email = c.Query("email")
	user.Username = c.Query("username")
	user.Password = c.Query("password")

	user, _, err := h.DB.AddUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("User with _id = %s added!", user.ID)})
	}
}

func GetAllUsers(c *gin.Context, h *Handler, origin *models.User) {
	users, err := h.DB.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context, h *Handler, origin *models.User) {
	userID := c.Param("user_id")
	user, err := h.DB.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(c *gin.Context, h *Handler, origin *models.User) {
	updates := make(map[string]interface{})
	userID := c.Param("user_id")
	for k, v := range c.Request.URL.Query() {
		if len(v) > 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
			return
		} else {
			updates[k] = v[0]
		}
	}

	_, err := h.DB.UpdateUser(userID, updates)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("User with _id = %s updated!", userID)})
	}
}

func DeleteUser(c *gin.Context, h *Handler, origin *models.User) {
	userID := c.Param("user_id")
	err := h.DB.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("User with _id = %s deleted!", userID)})
	}
}
