package web

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"net/http"
)

func ListServersController(c *gin.Context, db *gorm.DB) {
	var servers []models.Server
	db.Find(&servers)

	c.HTML(http.StatusOK, "servers", gin.H{
		"title":   "Server",
		"nav":     "servers",
		"servers": servers,
	})
}

func ShowServerController(c *gin.Context, db *gorm.DB, serverId string) {
	var server models.Server
	db.First(&server, serverId)

	if server.ID == 0 {
		c.Redirect(302, "/servers")
		return
	}

	c.HTML(http.StatusOK, "server", gin.H{
		"title":  "Server",
		"nav":    "servers",
		"server": server,
	})
}

func ShowCreateServerController(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, "servers-create", gin.H{
		"title": "Create server",
		"nav":   "servers",
	})
}

func CreateServerController(c *gin.Context, db *gorm.DB) {
	db.Create(&models.Server{
		Name:                 c.PostForm("name"),
		Address:              c.PostForm("address"),
		CfAccessClientId:     c.PostForm("cf_access_client_id"),
		CfAccessClientSecret: c.PostForm("cf_access_client_secret"),
	})

	c.Redirect(302, "/servers")
}
