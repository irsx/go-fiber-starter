package main

import (
	"fmt"
	"go-fiber-starter/app/middlewares"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/app/workers"
	"go-fiber-starter/configs"
	"go-fiber-starter/database/seeders"
	"go-fiber-starter/routes"
	"go-fiber-starter/utils"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {
	// load env first to avoid env error
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using env from machine")
	}

	// setup custom fiber config
	app := fiberConfig()

	// load custom logger
	utils.ZapLogger(os.Getenv("APP_ENV"))

	// load config
	configs.Setup()
	repository.DB = configs.DB // pass gorm.DB to repository

	// middlewares
	middlewares.Setup(app)

	// setup route
	routes.Setup(app)

	// listen custom args
	argsListener()

	// listen rabbitmq workers
	go workerListener()

	// listen app
	port := os.Getenv("APP_PORT")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func fiberConfig() *fiber.App {
	// setup app
	maxBody, _ := strconv.Atoi(os.Getenv("APP_MAX_BODY"))
	if maxBody == 0 {
		maxBody = 4
	}

	// custom fiber config
	app := fiber.New(fiber.Config{
		ReadBufferSize: 4096 * 4,                     // 2mb
		BodyLimit:      maxBody * 1024 * 1024,        // 4mb
		Views:          html.New("./views", ".html"), // Set the HTML template engine on views folder
	})

	// set static folder
	app.Static("/views", "./views")

	return app
}

func workerListener() {
	connection, channel, err := configs.RabbitMQConnection()
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()
	defer channel.Close()

	var listConsumers []interface{}
	listConsumers = append(listConsumers, new(workers.UserRegisterWorker))

	worker := workers.Worker{
		Channel:   channel,
		Consumers: listConsumers,
	}

	worker.Listen()
}

func argsListener() {
	for _, arg := range os.Args {
		if arg == "--rollback" {
			utils.Logger.Info("✅ down all migration")
			cmd := exec.Command("sql-migrate", "down", "-limit=0")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			utils.Logger.Info("✅ up all migration")
			cmd = exec.Command("sql-migrate", "up", "-limit=0")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

		if arg == "--seed" {
			// Run all seeders
			runner := seeders.All(configs.DB)
			if err := runner.Run(); err != nil {
				panic(err)
			}
		}
	}
}
