package middlewares

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *fiber.App) {
	// panic recovery
	app.Use(recover.New())

	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))

	// limit repeated requests
	maxRequest, _ := strconv.Atoi(os.Getenv("APP_MAX_REQUEST"))
	app.Use(Limit(maxRequest, 5))

	// http request logger
	app.Use(logger.New(logger.Config{
		Format:     `${time} ${locals:requestid} ${status} - ${method} ${url}` + "\n\n",
		TimeFormat: "2006/01/02 15:04:05",
	}))

	// debugger
	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if debug {
		app.Use(pprof.New())
	}
}

func Limit(maxRequest int, duration time.Duration) func(*fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        maxRequest,
		Expiration: duration * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error":   true,
				"message": "Too many requests",
			})
		},
	})
}

func MaxBodySize(sizeInMB int) fiber.Handler {
	sizeInMB = sizeInMB * 1024 * 1024
	return func(c *fiber.Ctx) error {
		if len(c.Body()) >= sizeInMB {
			// custom response here
			return fiber.ErrRequestEntityTooLarge
		}
		return c.Next()
	}
}
