package main

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/models"
	"github.com/romanzipp/wanderer/routes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

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

	app := &application.App{}
	app.DB = MakeDb()
	app.Router = MakeRouter(app.DB)
	app.Env = makeEnv()

	routes.InitApiRoutes(app)
	routes.InitWebRoutes(app)

	go MakeCheckScheduler(app)

	log.Fatal().Err(app.Router.Run())
}

func makeEnv() application.Env {
	lifetime, _ := strconv.ParseInt(os.Getenv("SESSION_LIFETIME"), 10, 64)

	return application.Env{
		Password:        os.Getenv("APP_PASSWORD"),
		SessionLifetime: int(lifetime),
	}
}

func MakeDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data/data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.Server{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.Template{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.TemplateVersion{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.Token{}); err != nil {
		panic(err)
	}

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

	return router
}

func MakeCheckScheduler(app *application.App) {
	for {
		var servers []models.Server
		app.DB.Find(&servers)

		for _, server := range servers {
			status, err := server.Check(app.DB)
			if err != nil {
				log.Warn().Msgf("[server check] %s error: %s", server.Name, err)
			} else {
				log.Debug().Msgf("[server check] %s = %d", server.Name, status)
			}
		}

		time.Sleep(60 * time.Second)
	}
}
