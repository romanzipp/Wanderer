package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitWebRoutes(router *gin.Engine) {
	// --------------------------------------------
	// authed routes

	authed := router.Group("/")
	authed.Use(AuthRequired())

	authed.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"title": "Main website",
		})
	})

	// --------------------------------------------
	// login

	router.GET("/auth", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth", gin.H{
			"title": "Authenticate",
		})
	})

	router.POST("/auth", func(c *gin.Context) {
		c.SetCookie("token", "ok", 3000, "/", "", false, true)
		c.Redirect(301, "/")
	})
}

func AuthRequired() gin.HandlerFunc {
	abort := func(c *gin.Context) {
		c.Redirect(301, "/auth")
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
