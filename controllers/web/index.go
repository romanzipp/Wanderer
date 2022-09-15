package web

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"net/http"
)

func WebIndexController(c *gin.Context, db *gorm.DB) {
	var templates []models.Template
	var servers []models.Server

	db.Model(&models.Template{}).Preload("Server").Preload("Versions").Find(&templates)
	db.Find(&servers)

	var selectedServer models.Server
	var selectedTemplate models.Template

	if c.Query("server") != "" {
		db.First(&selectedServer, c.Query("server"))
	}

	if c.Query("template") != "" {
		db.First(&selectedTemplate, c.Query("template"))
	}

	var versions []models.TemplateVersion

	if selectedTemplate.ID != 0 {
		db.Where("template_id = ?", selectedTemplate.ID).Find(&versions)
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"title":     "Home",
		"nav":       "index",
		"servers":   servers,
		"templates": templates,
		"server":    selectedServer,
		"template":  selectedTemplate,
		"versions":  versions,
	})
}
