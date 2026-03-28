package presenters

// import (
// 	"backend/models"

// 	"github.com/gofiber/fiber/v2"
// )

// func TicketSuccessResponse(data *models.Ticket) *fiber.Map {
// 	ticket := models.Ticket{
// 		TicketID:    data.TicketID,
// 		Subject:     data.Subject,
// 		Description: data.Description,
// 		RequestType: data.RequestType,
// 		Status:      data.Status,
// 		Priority:    data.Priority,
// 		AssignedTo:  data.AssignedTo,
// 		AssignedBy:  data.AssignedBy,
// 		CreatedBy:   data.CreatedBy,
// 		DueDate:     data.DueDate,
// 		Notes:       data.Notes,
// 		Attachments: data.Attachments,
// 		CreatedOn:   data.CreatedOn,
// 		ModifiedOn:  data.ModifiedOn,
// 	}

// 	return &fiber.Map{
// 		"status": "success",
// 		"data":   ticket,
// 		"error":  nil,
// 	}
// }

// // TicketsSuccessResponse is the list SuccessResponse that
// // will be passed in the response by Handler
// func TicketsSuccessResponse(data []*models.Ticket) *fiber.Map {
// 	return &fiber.Map{
// 		"status": "success",
// 		"data":   data,
// 		"error":  nil,
// 	}
// }

// // TicketErrorResponse is the ErrorResponse that
// // will be passed in the response by Handler
// func TicketErrorResponse(err error) *fiber.Map {
// 	return &fiber.Map{
// 		"status": "fail",
// 		"data":   "",
// 		"error":  err.Error(),
// 	}
// }
