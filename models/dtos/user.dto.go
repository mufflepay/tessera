package dtos

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

type RegisterDTO struct {
	FirstName       string   `json:"first_name" valid:"alphanum,required~first name cannot be empty"`
	LastName        string   `json:"last_name" valid:"alphanum,required~last name cannot be empty"`
	Email           string   `json:"email" valid:"email,required~email cannot be empty"`
	Password        string   `json:"password" valid:"length(8|100)~password min-length (8) is required,required~password cannot be empty"`
	ConfirmPassword string   `json:"confirm_password" valid:"length(8|100)~confirm password min-length (8) is required,required~confirm password cannot be empty"`
	PhotoURL        string   `json:"photo_url" valid:"url,optional"`
	Role            UserType `json:"role" valid:"in(admin|individual|corporate),required~role cannot be empty"`
}

func (r *RegisterDTO) IsRoleValid() bool {
	switch r.Role {
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

func (r *RegisterDTO) IsRoleInvalid() bool {
	return !r.IsRoleValid()
}

type LoginDTO struct {
	Email    string `json:"email"  valid:"email,required~email cannot be empty"`
	Password string `json:"password"  valid:"length(8|100)~password min length (8) is required,required~password cannot be empty"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" valid:"email,required~email cannot be empty"`
}

type ResendVerificationEmail struct {
	Email string `json:"email" valid:"email,required~email cannot be empty"`
}

type ResetPasswordDTO struct {
	Password        string `json:"password" valid:"length(8|100)~password min length (8) is required,required~password cannot be empty"`
	ConfirmPassword string `json:"confirm_password" valid:"length(8|100)~matching confirm password min length (8) is required,required~confirm password cannot be empty"`
}

type UserResponseDTO struct {
	UserID     uuid.UUID `json:"user_id,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Email      string    `json:"email,omitempty"`
	Role       string    `json:"role,omitempty"`
	Provider   string    `json:"provider"`
	PhotoURL   string    `json:"photo_url,omitempty"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
}

// func FilterUserRecord(user *models.User) UserResponseDTO {
// 	return UserResponseDTO{
// 		UserID:     user.UserID,
// 		FirstName:  user.FirstName,
// 		LastName:   user.LastName,
// 		Email:      user.Email,
// 		Role:       user.Role,
// 		Provider:   user.Provider,
// 		PhotoURL:   user.PhotoURL,
// 		CreatedOn:  user.CreatedOn,
// 		ModifiedOn: user.ModifiedOn,
// 	}
// }

// func FilterUserRecords(users []*models.User) []UserResponseDTO {
// 	var usersDTO []UserResponseDTO
// 	for _, user := range users {
// 		usersDTO = append(usersDTO, FilterUserRecord(user))
// 	}
// 	return usersDTO
// }

// type UserResponseAllDTO struct {
// 	UserID     uuid.UUID `json:"user_id,omitempty"`
// 	FirstName  string    `json:"first_name,omitempty"`
// 	LastName   string    `json:"last_name,omitempty"`
// 	Email      string    `json:"email,omitempty"`
// 	Password   string    `json:"password,omitempty"`
// 	Role       string    `json:"role,omitempty"`
// 	PhotoURL   string    `json:"photo_url,omitempty"`
// 	Verified   bool      `json:"verified,omitempty"`
// 	Provider   string    `json:"provider"`
// 	CreatedOn  time.Time `json:"created_on"`
// 	ModifiedOn time.Time `json:"modified_on"`
// }

// func FilterUserRecordWithPassword(user *models.User) UserResponseAllDTO {
// 	return UserResponseAllDTO{
// 		UserID:     user.UserID,
// 		FirstName:  user.FirstName,
// 		LastName:   user.LastName,
// 		Email:      user.Email,
// 		Password:   *user.Password,
// 		Role:       *user.Role,
// 		Provider:   user.Provider,
// 		PhotoURL:   *user.PhotoURL,
// 		CreatedOn:  user.CreatedOn,
// 		ModifiedOn: *user.ModifiedOn,
// 	}
// }
