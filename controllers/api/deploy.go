package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type DeployPayload struct {
	Server   int    `json:"server" binding:"required"`
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

func DeployController(c *gin.Context, db *gorm.DB) {
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
	db.First(&server, payload.Server)
	if server.ID == 0 {
		c.JSON(422, &ErrorResponse{"invalid server id"})
		return
	}

	// find template
	db.Where("nomad_job_id = ?", payload.Job).Preload("Server").Find(&template)
	if template.ID == 0 {
		c.JSON(422, &ErrorResponse{"invalid job id"})
		return
	}

	if template.ServerID != server.ID {
		c.JSON(422, &ErrorResponse{"template id and server id do not belong to the same server"})
		return
	}

	// get template version for provided selector
	db.Where("selector = ?", payload.Selector).Where("template_id = ?", template.ID).Find(&templateVersion)
	if templateVersion.ID == 0 {
		c.JSON(422, &ErrorResponse{"invalid template selector"})
		return
	}

	// set version string
	templateVersion.LastVersion = payload.Version

	// deploy
	err = template.Deploy(db, &templateVersion, payload.Version)
	if err != nil {
		log.Warn().Err(err)
		c.JSON(400, &ErrorResponse{err.Error()})
		return
	}

	c.JSON(200, SuccessResponse{"ok"})
}
