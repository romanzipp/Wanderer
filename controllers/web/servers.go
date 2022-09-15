package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

func ShowServerController(c *gin.Context, db *gorm.DB, serverID string) {
	var server models.Server
	db.First(&server, serverID)

	if server.ID == 0 {
		c.Redirect(302, "/servers")
		return
	}

	c.HTML(http.StatusOK, "server", gin.H{
		"title":   "Server",
		"nav":     "servers",
		"server":  server,
		"success": c.Query("success"),
	})
}

func UpdateServerController(c *gin.Context, db *gorm.DB, serverID string) {
	serverIDConv, _ := strconv.ParseInt(serverID, 10, 64)
	server := fillServer(c)
	server.ID = uint(serverIDConv)

	db.Save(server)

	c.Redirect(302, fmt.Sprintf("/servers/%s?success=Server+updated", serverID))
}

func ShowCreateServerController(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, "servers-create", gin.H{
		"title": "Create server",
		"nav":   "servers",
	})
}

func CreateServerController(c *gin.Context, db *gorm.DB) {
	server := fillServer(c)
	db.Create(&server)

	c.Redirect(302, "/servers?success=Server+created")
}

func fillServer(c *gin.Context) models.Server {
	return models.Server{
		Name:                 c.PostForm("name"),
		Address:              c.PostForm("address"),
		CfAccessClientId:     c.PostForm("cf_access_client_id"),
		CfAccessClientSecret: c.PostForm("cf_access_client_secret"),
	}
}
