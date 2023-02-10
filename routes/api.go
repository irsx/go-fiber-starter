package routes

import (
	"go-fiber-starter/app/controllers"
	"go-fiber-starter/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
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
	importGroup.Get("/products/events", importController.ImportProductsEvents)

	// PRODUCT ROUTES
	productGroup := apiRoute.Group("/product")
	productController := new(controllers.ProductController)
	productGroup.Get("/", productController.List)
	productGroup.Post("/", productController.Add)
	productGroup.Put("/:guid", productController.Update)
	productGroup.Delete("/:guid", productController.Delete)

	// UPLOAD ROUTES
	uploadController := new(controllers.UploadController)
	apiRoute.Post("/cdn", uploadController.CDN)
	apiRoute.Post("/callback", uploadController.Callback)
}
