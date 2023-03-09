package routes

import (
	"go-fiber-starter/app/controllers"
	"go-fiber-starter/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// PUBLIC ROUTES
	app.Get("/import-demo", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	apiRoute := app.Group("/api", middlewares.JwtMiddleware)

	// AUTH ROUTES
	auth := apiRoute.Group("/auth")
	authController := new(controllers.AuthController)
	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)

	// IMPORT ROUTES
	importGroup := apiRoute.Group("/import")
	importController := new(controllers.ImportController)
	importGroup.Get("/logs", importController.ImportLogAll)
	importGroup.Get("/logs/:guid", importController.ImportLogDetail)
	importGroup.Post("/products", importController.ImportProducts)
	importGroup.Get("/products/stream", importController.ImportProductsStream)

	// USER ROUTES
	userGroup := apiRoute.Group("/user")
	userController := new(controllers.UserController)
	userGroup.Get("/", userController.List)
	userGroup.Post("/", userController.Add)
	userGroup.Put("/:guid", userController.Update)
	userGroup.Delete("/:guid", userController.Delete)

	// PRODUCT ROUTES
	productGroup := apiRoute.Group("/product")
	productController := new(controllers.ProductController)
	productGroup.Get("/", productController.List)
	productGroup.Get("/:guid", productController.Detail)
	productGroup.Post("/", productController.Add)
	productGroup.Put("/:guid", productController.Update)
	productGroup.Delete("/:guid", productController.Delete)

	// NEWS ROUTES
	newsGroup := apiRoute.Group("/news")
	newsController := new(controllers.NewsController)
	newsGroup.Get("/", newsController.List)
	newsGroup.Get("/:guid", newsController.Detail)
	newsGroup.Post("/", newsController.Add)
	newsGroup.Put("/:guid", newsController.Update)
	newsGroup.Delete("/:guid", newsController.Delete)

	// UPLOAD ROUTES
	uploadController := new(controllers.UploadController)
	apiRoute.Post("/cdn", uploadController.CDN)
}
