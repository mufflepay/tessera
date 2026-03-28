package main

import (
	"backend/config/db"
	"backend/controllers"
	"backend/repos"
	"backend/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // "http://localhost:5173",
		AllowHeaders:     "Accept, Origin, Content-Type",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: false, // set true to pass the cookie header from client to backend
	}))
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(&fiber.Map{
				"status":  "fail",
				"message": "You have made too many requests in a single time-frame! Please wait a minute!",
			})
		},
	}))

	Connect(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("Welcome to the 😊Tessera Ticketing System Golang API😊")
	})

	go func() {
		if err := app.Listen(":8080"); err != nil && err != http.ErrServerClosed {
			log.Panicf("Failed to shutdown server: %s", err)
		}
	}()
	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	log.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}

	log.Println("Fiber was successful shutdown.")
}

func Connect(app *fiber.App) {
	// Connect to databases
	postgresDB, err := db.SetupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database. ", err.Error())
	}

	// Setup repositories and controllers
	userRepo := repos.NewUserRepository(postgresDB)
	ticketRepo := repos.NewTicketRepository(postgresDB)
	authController := controllers.NewAuthController(userRepo)
	authRoutes := routes.NewAuthRoutes(authController)
	userController := controllers.NewUserController(userRepo)
	userRoutes := routes.NewUserRoutes(userController)
	ticketController := controllers.NewTicketController(ticketRepo)
	ticketRoutes := routes.NewTicketRoutes(ticketController)

	// Register routes
	authRoutes.InitRoutes(app)
	userRoutes.InitRoutes(app)
	ticketRoutes.InitRoutes(app)
}
