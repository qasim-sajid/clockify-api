package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func SignUpUser(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &models.User{}

		err := c.ShouldBindJSON(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request format: %v", err).Error()})
			return
		}

		user, err = validateUser(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		addedUser, status, err := h.DB.AddUser(user)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, addedUser)
	}
}

func validateUser(user *models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, errors.New("name is missing")
	} else if user.Email == "" {
		return nil, errors.New("email is missing")
	} else if user.Username == "" {
		return nil, errors.New("username is missing")
	} else if user.Password == "" {
		return nil, errors.New("password is missing")
	}

	return user, nil
}
