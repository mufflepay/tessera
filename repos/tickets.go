package repos

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITicketRepository interface {
	CreateTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
	GetTickets(ctx context.Context) ([]*models.Ticket, error)
	GetTicketByID(ctx context.Context, ticketID uuid.UUID) (*models.Ticket, error)
	// UpdateTicket(ctx context.Context, ticketID uuid.UUID, ticketUpdate *models.Ticket) (*models.Ticket, error)
	//DeleteTicket(ctx context.Context, ticketID uuid.UUID) (uuid.UUID, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) ITicketRepository {
	return &ticketRepository{db: db}
}

// CreateTicket is a repository that helps to create tickets
func (r *ticketRepository) CreateTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {

	result := r.db.Create(&ticket)

	if result.Error != nil {
		return nil, result.Error
	}

	return ticket, nil
}

func (r *ticketRepository) GetTickets(ctx context.Context) ([]*models.Ticket, error) {
	var tickets []*models.Ticket

	result := r.db.Preload(clause.Associations).Find(&tickets)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return tickets, nil
}

func (r *ticketRepository) GetTicketByID(ctx context.Context, ticketID uuid.UUID) (*models.Ticket, error) {
	var ticket *models.Ticket

	result := r.db.Preload(clause.Associations).First(&ticket, ticketID)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return ticket, nil
}

//func (r *ticketRepository) UpdateTicket(ctx context.Context, ticketID uuid.UUID, ticketUpdate *models.Ticket) (*models.Ticket, error) {
//	var ticket models.Ticket
//
//	// Find the ticket first to get its associations
//	if err := r.db.Where("ticket_id = ?", ticketID).First(&ticket).Error; err != nil {
//		return nil, err
//	}
//
//	// Update the fields of the ticket with the values from ticketUpdate
//	if err := r.db.Model(&ticket).Updates(ticketUpdate).Error; err != nil {
//		return nil, err
//	}
//
//	// Now update the associations separately
//	if ticketUpdate.Status != nil {
//		// Assuming Status is a pointer to the new Status value
//		ticket.Status = ticketUpdate.Status
//	}
//
//	if ticketUpdate.Priority != nil {
//		// Assuming Priority is a pointer to the new Priority value
//		ticket.Priority = ticketUpdate.Priority
//	}
//
//	if len(ticketUpdate.Attachments) > 0 {
//		// Assuming Attachments is a slice of the new attachments
//		ticket.Attachments = ticketUpdate.Attachments
//	}
//
//	// Save the updated ticket with its associations
//	if err := r.db.Save(&ticket).Error; err != nil {
//		return nil, err
//	}
//
//	return &ticket, nil
//}
