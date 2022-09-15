package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ApiController(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, "api", gin.H{
		"title": "API",
		"nav":   "api",
	})
}
