package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/routes"
	"log"
)

func main() {
	router := MakeRouter()
	log.Fatal(router.Run())
}

func MakeRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")
	router.Static("/dist", "./dist")

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	routes.InitWebRoutes(router)

	return router
}
