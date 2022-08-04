package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddClient(c *gin.Context, h *Handler, origin *models.User) {
	client := &models.Client{}
	client.Name = c.Query("name")
	client.Address = c.Query("address")
	client.Note = c.Query("note")
	isArchived, err := strconv.ParseBool(c.Query("is_archived"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	client.IsArchived = isArchived

	client, _, err = h.DB.AddClient(client)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Client with _id = %s added!", client.ID)})
	}
}

func GetAllClients(c *gin.Context, h *Handler, origin *models.User) {
	clients, err := h.DB.GetAllClients()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clients)
}

func GetClient(c *gin.Context, h *Handler, origin *models.User) {
	clientID := c.Param("client_id")
	client, err := h.DB.GetClient(clientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, client)
	}
}

func UpdateClient(c *gin.Context, h *Handler, origin *models.User) {
	updates := make(map[string]interface{})
	clientID := c.Param("client_id")
	for k, v := range c.Request.URL.Query() {
		if len(v) > 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
			return
		} else {
			updates[k] = v[0]
		}
	}

	_, err := h.DB.UpdateClient(clientID, updates)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Client with _id = %s updated!", clientID)})
	}
}

func DeleteClient(c *gin.Context, h *Handler, origin *models.User) {
	clientID := c.Param("client_id")
	err := h.DB.DeleteClient(clientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Client with _id = %s deleted!", clientID)})
	}
}
