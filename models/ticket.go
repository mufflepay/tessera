package models

import (
	"time"

	"github.com/google/uuid"
)

// Ticket ==== ENTITIES ==== //
type Ticket struct {
	TicketID    uuid.UUID `json:"ticket_id" gorm:"column:ticket_id;primaryKey;not null;index;type:uuid;default:uuid_generate_v4()" validate:"omitempty,uuid4"`
	Subject     string    `json:"subject" gorm:"type:varchar(255);not null;index" validate:"required"`
	Description string    `json:"description" gorm:"type:text;not null" validate:"required"`
	//RequestType string              `json:"request_type" gorm:"type:varchar(50);not null;index" validate:"required"`
	//Status      *TicketStatus       `json:"status" gorm:"foreignKey:TicketID"`
	//Priority    *TicketPriority     `json:"priority" gorm:"foreignKey:TicketID"`
	//AssignedTo  string              `json:"assigned_to" gorm:"type:varchar(50);not null;index" validate:"required"`
	//AssignedBy  string              `json:"assigned_by" gorm:"type:varchar(50);not null;index" validate:"required"`
	//CreatedBy   string              `json:"created_by" gorm:"type:varchar(50);not null;index" validate:"required"`
	//DueDate     time.Time          `json:"due_date" gorm:"column:due_date;default:'0000-00-00'" validate:"required"`
	Attachments []*TicketAttachment `json:"attachments" gorm:"foreignKey:TicketID"`
	//IsArchived  bool                `json:"is_archived" gorm:"column:is_archived;not null;default:false" validate:"required"`
	//IsDeleted   bool                `json:"is_deleted" gorm:"column:is_deleted;not null;default:false" validate:"required"`
	//CreatedOn  time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	//ModifiedOn time.Time `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}

type TicketStatus struct {
	StatusID   uuid.UUID `json:"status_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4();"`
	TicketID   uuid.UUID `json:"ticket_id"`
	Status     string    `json:"status" gorm:"type:varchar(20);not null;index;default:Open" validate:"required"`
	ChangedBy  string    `json:"changed_by" gorm:"type:varchar(50);not null;index" validate:"required"`
	CreatedOn  time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn time.Time `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}

type TicketPriority struct {
	PriorityID uuid.UUID `json:"priority_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	TicketID   uuid.UUID `json:"ticket_id"`
	Priority   string    `json:"priority" gorm:"type:varchar(20);not null;index;default:High" validate:"required"`
	ChangedBy  string    `json:"changed_by" gorm:"type:varchar(50);not null;index" validate:"required"`
	CreatedOn  time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn time.Time `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}

type TicketAttachment struct {
	AttachmentID uuid.UUID `json:"attachment_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	TicketID     uuid.UUID `json:"ticket_id"`
	FileName     string    `json:"file_name" gorm:"column:file_name"`
	FileUrl      string    `json:"file_url" gorm:"column:file_url"`
	FileType     string    `json:"file_type" gorm:"column:file_type"`
	FileSize     int       `json:"file_size" gorm:"column:file_size"`
	CreatedOn    time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
}

type TicketStatusHistory struct {
	StatusID   uuid.UUID `json:"status_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4();"`
	TicketID   uuid.UUID `json:"ticket_id"`
	Status     string    `json:"status" gorm:"type:varchar(20);not null;index;default:Open" validate:"required"`
	ChangedBy  string    `json:"changed_by" gorm:"type:varchar(50);not null;index" validate:"required"`
	CreatedOn  time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn time.Time `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}

type TicketPriorityHistory struct {
	PriorityID uuid.UUID `json:"priority_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	TicketID   uuid.UUID `json:"ticket_id"`
	Priority   string    `json:"priority" gorm:"type:varchar(20);not null;index;default:High" validate:"required"`
	ChangedBy  string    `json:"changed_by" gorm:"type:varchar(50);not null;index" validate:"required"`
	CreatedOn  time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn time.Time `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}
