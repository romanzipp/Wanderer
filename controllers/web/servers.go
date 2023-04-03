package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/models"
	"net/http"
	"strconv"
)

func ListServersController(c *gin.Context, app *application.App) {
	var servers []models.Server
	app.DB.Find(&servers)

	c.HTML(http.StatusOK, "servers", gin.H{
		"title":      "Server",
		"nav":        "servers",
		"servers":    servers,
		"success":    c.Query("success"),
		"error":      c.Query("error"),
		"currentUrl": c.Request.URL.Path,
	})
}

func ShowServerController(c *gin.Context, app *application.App, serverID string) {
	var server models.Server
	app.DB.First(&server, serverID)

	if server.ID == 0 {
		c.Redirect(302, "/servers")
		return
	}

	c.HTML(http.StatusOK, "server", gin.H{
		"title":      "Server",
		"nav":        "servers",
		"server":     server,
		"success":    c.Query("success"),
		"currentUrl": c.Request.URL.Path,
	})
}

func UpdateServerController(c *gin.Context, app *application.App, serverID string) {
	serverIDConv, _ := strconv.ParseInt(serverID, 10, 64)
	server := fillServer(c)
	server.ID = uint(serverIDConv)

	app.DB.Save(server)

	c.Redirect(302, fmt.Sprintf("/servers/%s?success=Server+updated", serverID))
}

func DeleteServerController(c *gin.Context, app *application.App, serverID string) {
	var server models.Server
	app.DB.Where("id = ?", serverID).First(&server)

	if server.ID == 0 {
		c.Redirect(302, "/servers?error=Server+not+found")
		return
	}

	app.DB.Delete(&server)

	c.Redirect(302, "/servers?success=Server+deleted")
}

func ShowCreateServerController(c *gin.Context, app *application.App) {
	c.HTML(http.StatusOK, "servers-create", gin.H{
		"title":      "Create server",
		"nav":        "servers",
		"currentUrl": c.Request.URL.Path,
	})
}

func CreateServerController(c *gin.Context, app *application.App) {
	server := fillServer(c)
	app.DB.Create(&server)

	c.Redirect(302, "/servers?success=Server+created")
}

func fillServer(c *gin.Context) models.Server {
	return models.Server{
		Name:                 c.PostForm("name"),
		Address:              c.PostForm("address"),
		CfAccessClientId:     c.PostForm("cf_access_client_id"),
		CfAccessClientSecret: c.PostForm("cf_access_client_secret"),
		BasicAuthUser:        c.PostForm("basic_auth_user"),
		BasicAuthPassword:    c.PostForm("basic_auth_password"),
	}
}
