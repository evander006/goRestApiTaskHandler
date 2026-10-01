package routes

import (
	"restApiGo/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/user")
	//api.Get("/get-user", controllers.User)
	api.Post("/register", controllers.Register)
	//api.Get("/login", controllers.Login)
	//api.Get("/logout", controllers.Logout)

}
