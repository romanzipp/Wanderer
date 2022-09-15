package api

import (
	"github.com/gin-gonic/gin"
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

	c.JSON(200, SuccessResponse{"ok"})
}
