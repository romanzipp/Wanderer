package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/models"
	"github.com/rs/zerolog/log"
)

type DeployPayload struct {
	Server   int    `json:"server"`
	Job      string `json:"job" binding:"required"`      // nomad job id
	Selector string `json:"selector" binding:"required"` // version selector
	Version  string `json:"version" binding:"required"`  // version
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func DeployController(c *gin.Context, app *application.App) {
	var payload DeployPayload
	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(422, &ErrorResponse{"invalid payload"})
		log.Warn().Err(err)
		return
	}

	var server models.Server
	var template models.Template
	var templateVersion models.TemplateVersion

	// find server
	if payload.Server == 0 {
		var count int64
		app.DB.Model(models.Server{}).Count(&count)
		if count != 1 {
			c.JSON(422, &ErrorResponse{"server parameter required since there is more than 1 server"})
			return
		}
		app.DB.First(&server)
	} else {
		app.DB.First(&server, payload.Server)
	}

	if server.ID == 0 {
		c.JSON(422, &ErrorResponse{"invalid server id"})
		return
	}

	// find template
	app.DB.Where("nomad_job_id = ?", payload.Job).Preload("Server").Find(&template)
	if template.ID == 0 {
		c.JSON(422, &ErrorResponse{fmt.Sprintf("invalid job id: %s", payload.Job)})
		return
	}

	if template.ServerID != server.ID {
		c.JSON(422, &ErrorResponse{"template id and server id do not belong to the same server"})
		return
	}

	// get template version for provided selector
	app.DB.Where("selector = ?", payload.Selector).Where("template_id = ?", template.ID).Find(&templateVersion)
	if templateVersion.ID == 0 {
		c.JSON(422, &ErrorResponse{"invalid template selector"})
		return
	}

	// set version string
	templateVersion.LastVersion = payload.Version

	// deploy
	err = template.Deploy(app.DB, &templateVersion, payload.Version)
	if err != nil {
		log.Warn().Err(err)
		c.JSON(400, &ErrorResponse{err.Error()})
		return
	}

	c.JSON(200, SuccessResponse{"ok"})
}
