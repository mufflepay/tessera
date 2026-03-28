package models

import (
	"github.com/google/uuid"
	"time"
)

type TicketNote struct {
	NoteID     uuid.UUID `json:"note_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	TicketID   uuid.UUID `json:"ticket_id"`
	Title      string    `json:"title" gorm:"column:title"`
	Content    string    `json:"content" gorm:"column:content"`
	Author     string    `json:"author" gorm:"column:author"`
	CreatedOn  time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn time.Time `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}
