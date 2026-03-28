package routes

import (
	"tessera/controllers"

	"github.com/gofiber/fiber/v2"
)

type ITicketRoutes interface {
	InitRoutes(app *fiber.App)
}

type ticketRoutes struct {
	ticketController controllers.ITicketController
}

func NewTicketRoutes(ticketController controllers.ITicketController) ITicketRoutes {
	return &ticketRoutes{ticketController}
}

func (t *ticketRoutes) InitRoutes(app *fiber.App) {
	ticketRoute := app.Group("/api/v1/tickets")

	// Declare routing endpoints for general routes.
	ticketRoute.Get("", t.ticketController.GetTickets)
	ticketRoute.Post("", t.ticketController.CreateTicket)

	// Declare routing endpoints for specific routes.
	ticketRoute.Get("/:id", t.ticketController.GetTicketByID)
	//ticketRoute.Put("/:id", t.ticketController.UpdateTicket)
	//ticketRoute.Delete("/:id", t.ticketController.DeleteTicket)

}
