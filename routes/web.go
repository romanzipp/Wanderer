package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/controllers/web"
	"gorm.io/gorm"
)

func InitWebRoutes(router *gin.Engine, db *gorm.DB) {
	// --------------------------------------------
	// authed routes

	authed := router.Group("/")
	authed.Use(AuthRequired())

	authed.GET("/", func(c *gin.Context) {
		web.WebIndexController(c, db)
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

	// --------------------------------------------
	// login

	router.GET("/auth", func(c *gin.Context) {
		web.ShowAuthFormController(c, db)
	})

	router.POST("/auth", func(c *gin.Context) {
		web.SubmitAuthController(c, db)
	})
}

func AuthRequired() gin.HandlerFunc {
	abort := func(c *gin.Context) {
		c.Redirect(302, "/auth")
	}

	return func(c *gin.Context) {
		token, _ := c.Request.Cookie("token")

		if len(token.String()) == 0 {
			abort(c)
			return
		}

		fmt.Println(token)
	}
}
