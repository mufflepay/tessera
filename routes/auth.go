package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
)

type AuthRoutes interface {
	InitRoutes(app *fiber.App)
}

type authRoutes struct {
	authController controllers.IAuthController
}

func NewAuthRoutes(authController controllers.IAuthController) AuthRoutes {
	return &authRoutes{authController}
}

func (a *authRoutes) InitRoutes(app *fiber.App) {
	authRoute := app.Group("/api/v1/auth")

	// Declare routing endpoints for general routes.
	authRoute.Post("/register", a.authController.RegisterUser)
	authRoute.Post("/login", a.authController.LoginUser)
	authRoute.Get("/logout", a.authController.LogoutUser)

	authRoute.Get("/verify-email/:verificationToken", a.authController.VerifyEmail)
	authRoute.Post("/resend-verification-email", a.authController.ResendVerificationEmail)
	authRoute.Get("/forgot-password", a.authController.ForgotPassword)
	authRoute.Post("/reset-password", a.authController.ResetPassword)
	authRoute.Get("/reset-password/:passResetToken", a.authController.ValidatePasswordResetToken)
	authRoute.Get("/revoke-session/:sessionID", a.authController.RevokeSession)
}
