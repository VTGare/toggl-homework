package main

import (
	"log"
	"os"

	"github.com/VTGare/toggl-homework/controllers"
	"github.com/VTGare/toggl-homework/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = os.Getenv("PORT")
)

func initFiber(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	controllers.RegisterQuestion(app)
}

func main() {
	app := fiber.New()

	if port == "" {
		port = "3000"
	}

	err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DBConn.Close()

	initFiber(app)

	if err := app.Listen(":" + port); err != nil {
		database.DBConn.Close()
		log.Fatal(err)
	}
}
