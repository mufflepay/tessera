package dtos

// import (
// 	"backend/models"
// 	"time"

// 	"github.com/google/uuid"
// )

// // ==== DTOs ==== //

// type CreateTicketDTO struct {
// 	Subject     string                   `json:"subject"`
// 	Description string                   `json:"description"`
// 	RequestType string                   `json:"request_type"`
// 	Status      []*models.TicketStatus   `json:"status"`
// 	Priority    []*models.TicketPriority `json:"priority"`
// 	AssignedTo  string                   `json:"assigned_to"`
// 	AssignedBy  string                   `json:"assigned_by"`
// 	CreatedBy   string                   `json:"created_by"`
// 	DueDate     time.Time                `json:"due_date"`
// 	Notes       []*models.Note           `json:"notes"`
// 	Attachment  []*models.Attachment     `json:"attachment"`
// 	CreatedOn   time.Time                `json:"created_on"`
// 	ModifiedOn  time.Time                `json:"modified_on"`
// }

// type TicketDTO struct {
// 	TicketID    uuid.UUID                `json:"ticket_id"`
// 	Subject     string                   `json:"subject"`
// 	Description string                   `json:"description"`
// 	RequestType string                   `json:"request_type"`
// 	Status      []*models.TicketStatus   `json:"status"`
// 	Priority    []*models.TicketPriority `json:"priority"`
// 	AssignedTo  string                   `json:"assigned_to"`
// 	AssignedBy  string                   `json:"assigned_by"`
// 	CreatedBy   string                   `json:"created_by"`
// 	DueDate     time.Time                `json:"due_date"`
// 	Notes       []*models.Note           `json:"notes"`
// 	Attachments []*models.Attachment     `json:"attachment"`
// 	CreatedOn   time.Time                `json:"created_on"`
// 	ModifiedOn  time.Time                `json:"modified_on"`
// }

// type FetchTicketDTO struct {
// 	TicketID uuid.UUID `json:"ticket_id"`
// }
