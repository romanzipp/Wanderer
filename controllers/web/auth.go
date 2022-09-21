package web

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ShowAuthFormController(c *gin.Context, app *application.App) {
	c.HTML(http.StatusOK, "auth", gin.H{
		"title": "Authenticate",
		"nav":   "index",
		"error": c.Query("error"),
	})
}

func SubmitAuthController(c *gin.Context, app *application.App) {
	password := c.PostForm("password")
	if password != app.Env.Password {
		log.Info().Msg("failed auth attempt")
		c.Redirect(302, "/auth?error=Invalid+password")
		return
	}

	c.SetCookie("token", "ok", 3000, "/", "", false, true)
	c.Redirect(302, "/")
}
