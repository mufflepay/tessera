package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPasswordReset struct {
	ID             uuid.UUID  `json:"id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	UserID         uuid.UUID  `json:"user_id" gorm:"not null;index;type:uuid"`
	PassResetToken string     `json:"pass_reset_token" gorm:"type:varchar(250);not null"`
	ExpiresOn      *time.Time `json:"expires_on" gorm:"not null"`
	PassResetOn    *time.Time `json:"pass_reset_on" gorm:"null"`
	CreatedOn      time.Time  `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn     time.Time  `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}
