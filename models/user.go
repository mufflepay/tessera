package models

import (
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	Admin      UserType = "admin"
	Individual UserType = "individual"
	Corporate  UserType = "corporate"
)

type User struct {
	UserID        uuid.UUID          `json:"user_id" gorm:"primaryKey;not null;index;type:uuid;default:uuid_generate_v4()"`
	FirstName     string             `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName      string             `json:"last_name" gorm:"type:varchar(100);not null"`
	Email         string             `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password      string             `json:"password" gorm:"type:varchar(100);not null"`
	Role          UserType           `json:"role" gorm:"check:role in('admin','individual','corporate');default:'individual';not null"`
	Provider      string             `json:"provider" gorm:"type:varchar(50);default:'local';not null"`
	PhotoURL      string             `json:"photo_url" gorm:"not null;default:'default.png'"`
	LastAccessAt  *time.Time         `json:"last_access_at" gorm:"null"`
	Verification  *UserVerification  `json:"verification" gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PasswordReset *UserPasswordReset `json:"password_reset" gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedOn     time.Time          `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn    time.Time          `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}

func (u *User) IsRoleValid() bool {
	switch u.Role {
	case UserType(Corporate):
		fallthrough
	case UserType(Admin):
		fallthrough
	case UserType(Individual):
		return true
	default:
		return false
	}
}

func (u *User) IsRoleInvalid() bool {
	return !u.IsRoleValid()
}
