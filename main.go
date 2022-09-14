package main

import (
	"fmt"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"github.com/romanzipp/wanderer/routes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	db := MakeDb()
	router := MakeRouter(db)

	go MakeCheckScheduler(db)

	log.Fatal(router.Run())
}

func MakeDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Job{})
	db.AutoMigrate(&models.Server{})
	db.AutoMigrate(&models.Template{})

	return db
}

func MakeRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.HTMLRender = ginview.Default()

	router.Static("/dist", "./dist")
	router.Static("/assets", "./static")

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	routes.InitWebRoutes(router, db)

	return router
}

func MakeCheckScheduler(db *gorm.DB) {
	for {
		var servers []models.Server
		db.Find(&servers)

		for _, server := range servers {
			status, err := server.Check(db)
			if err != nil {
				fmt.Printf("[server check] err %s\n", err)
			} else {
				fmt.Printf("[server check] %s = %d \n", server.Name, status)
			}
		}

		time.Sleep(60 * time.Second)
	}
}
