package main

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"github.com/romanzipp/wanderer/routes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	db := MakeDb()
	router := MakeRouter(db)

	go MakeCheckScheduler(db)

	log.Fatal().Err(router.Run())
}

func MakeDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Server{})
	db.AutoMigrate(&models.Template{})
	db.AutoMigrate(&models.TemplateVersion{})

	return db
}

func MakeRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.HTMLRender = ginview.Default()

	router.Static("/dist", "./dist")
	router.Static("/assets", "./static")

	router.Use(logger.SetLogger())

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
				log.Debug().Msgf("[server check] err %s", err)
			} else {
				log.Debug().Msgf("[server check] %s = %d", server.Name, status)
			}
		}

		time.Sleep(60 * time.Second)
	}
}
