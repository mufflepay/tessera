package models

import (
	"time"

	"github.com/google/uuid"
)

type UserVerification struct {
	ID                uuid.UUID  `json:"id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	UserID            uuid.UUID  `json:"user_id" gorm:"not null;index;type:uuid"`
	VerificationToken string     `json:"verification_token" gorm:"type:varchar(250);not null"`
	ExpiresOn         *time.Time `json:"expires_on" gorm:"not null"`
	IsVerified        bool       `json:"is_verified" gorm:"not null;default:false"`
	VerifiedOn        *time.Time `json:"verified_on" gorm:"null"`
	CreatedOn         time.Time  `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn        time.Time  `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}
