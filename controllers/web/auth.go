package web

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
	"net/http"
)

func ShowAuthFormController(c *gin.Context, app *application.App) {
	c.HTML(http.StatusOK, "auth", gin.H{
		"title": "Authenticate",
		"nav":   "index",
	})
}

func SubmitAuthController(c *gin.Context, app *application.App) {
	c.SetCookie("token", "ok", 3000, "/", "", false, true)
	c.Redirect(302, "/")
}
