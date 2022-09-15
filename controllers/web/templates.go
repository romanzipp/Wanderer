package web

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"net/http"
)

func ListTemplatesController(c *gin.Context, db *gorm.DB) {
	var templates []models.Template
	db.Find(&templates)

	c.HTML(http.StatusOK, "templates", gin.H{
		"title":     "Templates",
		"nav":       "templates",
		"templates": templates,
	})
}
