package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/controllers/web"
	"gorm.io/gorm"
)

func InitWebRoutes(router *gin.Engine, db *gorm.DB) {
	// --------------------------------------------
	// authed routes

	authed := router.Group("/")
	authed.Use(WebTokenAuth())

	authed.GET("/", func(c *gin.Context) {
		web.IndexController(c, db)
	})

	// servers

	authed.GET("/servers", func(c *gin.Context) {
		web.ListServersController(c, db)
	})

	authed.GET("/servers/:id", func(c *gin.Context) {
		web.ShowServerController(c, db, c.Param("id"))
	})

	authed.GET("/servers/create", func(c *gin.Context) {
		web.ShowCreateServerController(c, db)
	})

	authed.POST("/servers", func(c *gin.Context) {
		web.CreateServerController(c, db)
	})

	// templates

	authed.GET("/templates", func(c *gin.Context) {
		web.ListTemplatesController(c, db)
	})

	authed.GET("/templates/:id", func(c *gin.Context) {
		web.ShowTemplateController(c, db, c.Param("id"))
	})

	authed.POST("/templates/:id/versions", func(c *gin.Context) {
		web.CreateVersionController(c, db, c.Param("id"))
	})

	authed.GET("/templates/create", func(c *gin.Context) {
		web.ShowCreateTemplateController(c, db)
	})

	authed.POST("/templates", func(c *gin.Context) {
		web.CreateTemplateController(c, db)
	})

	// api

	authed.GET("/tokens", func(c *gin.Context) {
		web.ApiController(c, db)
	})

	authed.POST("/tokens", func(c *gin.Context) {
		web.IssueApiTokenController(c, db)
	})

	// --------------------------------------------
	// login

	router.GET("/auth", func(c *gin.Context) {
		web.ShowAuthFormController(c, db)
	})

	router.POST("/auth", func(c *gin.Context) {
		web.SubmitAuthController(c, db)
	})
}

func WebTokenAuth() gin.HandlerFunc {
	abort := func(c *gin.Context) {
		c.Redirect(302, "/auth")
	}

	return func(c *gin.Context) {
		token, _ := c.Request.Cookie("token")

		if len(token.String()) == 0 {
			abort(c)
			return
		}
	}
}
