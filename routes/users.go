package routes

import (
	"tessera/controllers"
	"tessera/middlewares"

	"github.com/gofiber/fiber/v2"
)

type UserRoutes interface {
	InitRoutes(app *fiber.App)
}

type userRoutes struct {
	userController controllers.UserController
}

func NewUserRoutes(userController controllers.UserController) UserRoutes {
	return &userRoutes{userController}
}

func (u *userRoutes) InitRoutes(app *fiber.App) {
	userRoute := app.Group("/api/v1/users")

	// Unrestricted routes to authenticated users
	userRoute.Get("/me", middlewares.AuthorizeAll, u.userController.GetMe)

	// Admin routes
	userRoute.Get("", middlewares.AuthorizeAdmin, u.userController.GetUsers)

	// Admin and Corporate routes
	userRoute.Get("/:id", middlewares.AuthorizeAdminOrCorporate, u.userController.GetUser)
}
