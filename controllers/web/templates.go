package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/nomad/api"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/models"
	"net/http"
	"strconv"
)

func ListTemplatesController(c *gin.Context, app *application.App) {
	var templates []models.Template
	app.DB.Model(&models.Template{}).Preload("Server").Preload("Versions").Find(&templates)

	c.HTML(http.StatusOK, "templates", gin.H{
		"title":     "Templates",
		"nav":       "templates",
		"templates": templates,
		"error":     c.Query("error"),
		"success":   c.Query("success"),
	})
}

func ShowTemplateController(c *gin.Context, app *application.App, templateID string) {
	var template models.Template
	app.DB.Preload("Versions").First(&template, templateID)

	if template.ID == 0 {
		c.Redirect(302, "/templates")
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	endpoint := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	c.HTML(http.StatusOK, "template", gin.H{
		"title":    "Template",
		"nav":      "templates",
		"template": template,
		"success":  c.Query("success"),
		"endpoint": endpoint,
	})
}

func UpdateTemplateController(c *gin.Context, app *application.App, templateID string) {
	serverID, _ := strconv.ParseInt(c.PostForm("server"), 10, 64)
	templateIDConv, _ := strconv.ParseInt(templateID, 10, 64)
	template := fillTemplate(c, uint(serverID))
	template.ID = uint(templateIDConv)

	app.DB.Save(template)

	c.Redirect(302, fmt.Sprintf("/templates/%s?success=Template+updated", templateID))
}

func DeleteTemplateController(c *gin.Context, app *application.App, templateID string) {
	var template models.Template
	app.DB.Where("id = ?", templateID).First(&template)

	if template.ID == 0 {
		c.Redirect(302, "/templates?error=Template+not+found")
		return
	}

	app.DB.Delete(&template)

	c.Redirect(302, "/templates?success=Template+deleted")
}

func ShowCreateTemplateController(c *gin.Context, app *application.App) {
	var servers []models.Server
	app.DB.Find(&servers)

	var selectedServer models.Server
	if c.Query("server") != "" {
		app.DB.First(&selectedServer, c.Query("server"))
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

func CreateTemplateController(c *gin.Context, app *application.App) {
	serverID, _ := strconv.ParseInt(c.PostForm("server"), 10, 64)
	template := fillTemplate(c, uint(serverID))

	app.DB.Create(&template)

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
