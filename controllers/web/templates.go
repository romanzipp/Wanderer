package web

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/nomad/api"
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
	db.Create(&models.Template{
		Name:       c.PostForm("name"),
		NomadJobID: c.PostForm("job"),
		Content:    c.PostForm("content"),
	})

	c.Redirect(302, "/servers")
}
