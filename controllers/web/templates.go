package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/nomad/api"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ListTemplatesController(c *gin.Context, db *gorm.DB) {
	var templates []models.Template
	db.Model(&models.Template{}).Preload("Server").Preload("Versions").Find(&templates)

	c.HTML(http.StatusOK, "templates", gin.H{
		"title":     "Templates",
		"nav":       "templates",
		"templates": templates,
	})
}

func ShowTemplateController(c *gin.Context, db *gorm.DB, templateID string) {
	var template models.Template
	db.Preload("Versions").First(&template, templateID)

	if template.ID == 0 {
		c.Redirect(302, "/templates")
		return
	}

	c.HTML(http.StatusOK, "template", gin.H{
		"title":    "Template",
		"nav":      "templates",
		"template": template,
		"success":  c.Query("success"),
	})
}

func UpdateTemplateController(c *gin.Context, db *gorm.DB, templateID string) {
	serverID, _ := strconv.ParseInt(c.PostForm("server"), 10, 64)
	templateIDConv, _ := strconv.ParseInt(templateID, 10, 64)
	template := fillTemplate(c, uint(serverID))
	template.ID = uint(templateIDConv)

	db.Save(template)

	c.Redirect(302, fmt.Sprintf("/templates/%s?success=Template+updated", templateID))
}

func ShowCreateTemplateController(c *gin.Context, db *gorm.DB) {
	var servers []models.Server
	db.Find(&servers)

	var selectedServer models.Server
	if c.Query("server") != "" {
		db.First(&selectedServer, c.Query("server"))
	}

	var nomadJobs []*api.JobListStub
	if selectedServer.ID != 0 {
		client, err := selectedServer.NewNomadClient()
		if err != nil {
			c.Redirect(302, "/templates/create?error=nomad-client")
			return
		}

		nomadJobs, _, err = client.Jobs().List(&api.QueryOptions{})
		if err != nil {
			c.Redirect(302, "/templates/create?error=nomad-jobs")
			return
		}
	}

	c.HTML(http.StatusOK, "templates-create", gin.H{
		"title":   "Create template",
		"nav":     "templates",
		"servers": servers,
		"server":  selectedServer,
		"jobs":    nomadJobs,
	})
}

func CreateTemplateController(c *gin.Context, db *gorm.DB) {
	serverID, _ := strconv.ParseInt(c.PostForm("server"), 10, 64)
	template := fillTemplate(c, uint(serverID))

	db.Create(&template)

	c.Redirect(302, "/templates?success=Template+created")
}

func fillTemplate(c *gin.Context, serverID uint) models.Template {
	return models.Template{
		Name:       c.PostForm("name"),
		NomadJobID: c.PostForm("job"),
		Content:    c.PostForm("content"),
		ServerID:   serverID,
	}
}
