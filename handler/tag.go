package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/models"
)

func AddTag(c *gin.Context, h *Handler, origin *models.User) {
	tag := &models.Tag{}
	tag.Name = c.Query("name")

	tag, _, err := h.DB.AddTag(tag)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Tag with _id = %s added!", tag.ID)})
	}
}

func GetAllTags(c *gin.Context, h *Handler, origin *models.User) {
	tags, err := h.DB.GetAllTags()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, tags)
}

func GetTag(c *gin.Context, h *Handler, origin *models.User) {
	tagID := c.Param("tag_id")
	tag, err := h.DB.GetTag(tagID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, tag)
	}
}

func UpdateTag(c *gin.Context, h *Handler, origin *models.User) {
	updates := make(map[string]interface{})
	tagID := c.Param("tag_id")
	for k, v := range c.Request.URL.Query() {
		if len(v) > 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Duplicate parameter found!"})
		} else {
			updates[k] = v[0]
		}
	}

	_, err := h.DB.UpdateTag(tagID, updates)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Tag with _id = %s updated!", tagID)})
	}
}

func DeleteTag(c *gin.Context, h *Handler, origin *models.User) {
	tagID := c.Param("tag_id")
	err := h.DB.DeleteTag(tagID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Tag with _id = %s deleted!", tagID)})
	}
}
