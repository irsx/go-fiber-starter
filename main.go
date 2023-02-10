package main

import (
	"fmt"
	"go-fiber-starter/app/middlewares"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/app/workers"
	"go-fiber-starter/configs"
	"go-fiber-starter/routes"
	"go-fiber-starter/utils"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using env from machine")
	}

	// setup app
	maxBody, _ := strconv.Atoi(os.Getenv("APP_MAX_BODY"))
	if maxBody == 0 {
		maxBody = 4
	}

	app := fiber.New(fiber.Config{
		ReadBufferSize: 4096 * 4,              // 2mb
		BodyLimit:      maxBody * 1024 * 1024, // 4mb
	})

	// load utils
	utils.ZapLogger(os.Getenv("APP_ENV"))

	// load config
	configs.Setup()

	// middlewares
	middlewares.Setup(app)

	// setup route
	routes.Setup(app)

	// pass gorm.DB to repository
	repository.DB = configs.DB

	// listen rabbitmq workers
	go workerListener()

	// listen app
	port := os.Getenv("APP_PORT")
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}

func workerListener() {
	connection, channel, err := configs.RabbitMQ()
	if err != nil {
		panic(err)
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
