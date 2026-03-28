package models

import (
	"time"

	"github.com/google/uuid"
)

type Provider struct {
	ProviderID   uuid.UUID `json:"provider_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	ProviderType string    `json:"provider_type" gorm:"type:varchar(50);not null"`
	PhotoURL     string    `json:"photo_url" gorm:"not null;default:'default.png'"`
	CreatedOn    time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
}
