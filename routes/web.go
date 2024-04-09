package routes

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/controllers/web"
)

func InitWebRoutes(app *application.App) {
	// --------------------------------------------
	// authed routes

	authed := app.Router.Group("/")
	authed.Use(WebTokenAuth(app))

	authed.GET("/", func(c *gin.Context) {
		web.IndexController(c, app.DB)
	})

	// servers

	authed.GET("/servers", func(c *gin.Context) {
		web.ListServersController(c, app)
	})

	authed.GET("/servers/:id", func(c *gin.Context) {
		web.ShowServerController(c, app, c.Param("id"))
	})

	authed.POST("/servers/:id", func(c *gin.Context) {
		web.UpdateServerController(c, app, c.Param("id"))
	})

	authed.POST("/servers/:id/delete", func(c *gin.Context) {
		web.DeleteServerController(c, app, c.Param("id"))
	})

	authed.GET("/servers/create", func(c *gin.Context) {
		web.ShowCreateServerController(c, app)
	})

	authed.POST("/servers", func(c *gin.Context) {
		web.CreateServerController(c, app)
	})

	// templates

	authed.GET("/templates", func(c *gin.Context) {
		web.ListTemplatesController(c, app)
	})

	authed.GET("/templates/:templateID", func(c *gin.Context) {
		web.ShowTemplateController(c, app, c.Param("templateID"))
	})

	authed.POST("/templates/:templateID", func(c *gin.Context) {
		web.UpdateTemplateController(c, app, c.Param("templateID"))
	})

	authed.POST("/templates/:templateID/lock", func(c *gin.Context) {
		web.LockTemplateController(c, app, c.Param("templateID"))
	})

	authed.POST("/templates/:templateID/delete", func(c *gin.Context) {
		web.DeleteTemplateController(c, app, c.Param("templateID"))
	})

	authed.POST("/templates/:templateID/versions", func(c *gin.Context) {
		web.CreateVersionController(c, app, c.Param("templateID"))
	})

	authed.POST("/templates/:templateID/redeploy", func(c *gin.Context) {
		web.RedeployController(c, app, c.Param("templateID"))
	})

	authed.POST("/versions/:versionID", func(c *gin.Context) {
		web.DeleteVersionController(c, app, c.Param("versionID"))
	})

	authed.GET("/templates/create", func(c *gin.Context) {
		web.ShowCreateTemplateController(c, app)
	})

	authed.POST("/templates", func(c *gin.Context) {
		web.CreateTemplateController(c, app)
	})

	// api

	authed.GET("/tokens", func(c *gin.Context) {
		web.ApiController(c, app)
	})

	authed.POST("/tokens", func(c *gin.Context) {
		web.IssueApiTokenController(c, app)
	})

	authed.POST("/tokens/:tokenID/delete", func(c *gin.Context) {
		web.DeleteApiTokenController(c, app, c.Param("tokenID"))
	})

	// --------------------------------------------
	// login

	app.Router.GET("/auth", func(c *gin.Context) {
		web.ShowAuthFormController(c, app)
	})

	app.Router.POST("/auth", func(c *gin.Context) {
		web.SubmitAuthController(c, app)
	})
}

func WebTokenAuth(app *application.App) gin.HandlerFunc {
	abort := func(c *gin.Context) {
		c.Redirect(302, "/auth")
	}

	return func(c *gin.Context) {
		token, _ := c.Request.Cookie("token")

		if len(token.String()) == 0 {
			abort(c)
			return
		}

		needle := sha256.Sum256([]byte(app.Env.Password))
		hash := fmt.Sprintf("%x", needle[:])

		if token.Value != hash {
			abort(c)
			return
		}

		c.Next()
	}
}
