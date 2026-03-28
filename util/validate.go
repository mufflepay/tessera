package util

import (
	"github.com/asaskevich/govalidator"
)

var (
	ErrInvalidCredentials        = "invalid credentials or account doesn't exists"
	ErrInvalidEmail              = "invalid email"
	ErrUserNotFound              = "user not found"
	ErrEmailAlreadyExists        = "email already exists"
	ErrEmailNotFound             = "email not found"
	ErrEmailNotSent              = "email not sent, please try again"
	ErrEmptyEmail                = "email cannot be empty"
	ErrEmptyPassword             = "password cannot be empty"
	ErrInvalidAuthTokens         = "invalid auth token"
	ErrUnauthorized              = "unauthorized"
	ErrPasswordNotMatch          = "passwords do not match"
	ErrInvalidUrl                = "invalid url"
	ErrEmailNotVerified          = "email not verified"
	ErrVerificationTokenRequired = "verification token required"
	ErrVerificationTokenExpired  = "verification token expired or invalid, please request a new one"
	ErrEmailAlreadyVerified      = "email already verified"
	ErrInvalidToken              = "invalid or expired token"
	ErrInvalidRole               = "invalid role"

	ValidTicketKeys = map[string]bool{
		"ticket_id":    true,
		"subject":      true,
		"description":  true,
		"request_type": true,
		"status":       true,
		"priority":     true,
		"assigned_to":  true,
		"assigned_by":  true,
		"created_by":   true,
		"due_date":     true,
		"notes":        true,
		"attachments":  true,
		"created_on":   true,
		"modified_on":  true,
	}
	ExcludeTicketKeys = map[string]bool{
		"created_by":  true,
		"created_on":  true,
		"due_date":    true,
		"status":      true,
		"priority":    true,
		"notes":       true,
		"attachments": true,
	}
	ValidNotesKeys = map[string]bool{
		"notes":      true,
		"note_id":    true,
		"ticket_id":  true,
		"content":    true,
		"author":     true,
		"created_on": true,
	}
	ExcludeNotesKeys = map[string]bool{
		"note_id":     true,
		"ticket_id":   true,
		"author":      true,
		"created_on":  true,
		"created_by":  true,
		"due_date":    true,
		"status":      true,
		"priority":    true,
		"attachments": true,
	}
	ValidAttachmentsKeys = map[string]bool{
		"attachments":   true,
		"attachment_id": true,
		"ticket_id":     true,
		"file_name":     true,
		"file_type":     true,
		"file_size":     true,
		"file_path":     true,
		"created_on":    true,
	}
	ExcludeAttachmentsKeys = map[string]bool{
		"attachment_id": true,
		"ticket_id":     true,
		"created_on":    true,
		"created_by":    true,
		"due_date":      true,
		"status":        true,
		"priority":      true,
		"notes":         true,
	}
	ValidStatusesKeys = map[string]bool{
		"status":     true,
		"status_id":  true,
		"ticket_id":  true,
		"old_status": true,
		"new_status": true,
		"changed_by": true,
	}
	ExcludeStatusesKeys = map[string]bool{
		"status_id":   true,
		"ticket_id":   true,
		"old_status":  true,
		"created_on":  true,
		"modified_on": true,
		"created_by":  true,
		"due_date":    true,
		"priority":    true,
		"notes":       true,
		"attachments": true,
	}
	ValidPrioritiesKeys = map[string]bool{
		"priority":     true,
		"priority_id":  true,
		"ticket_id":    true,
		"old_priority": true,
		"new_priority": true,
		"changed_by":   true,
	}
	ExcludePrioritiesKeys = map[string]bool{
		"priority_id":  true,
		"ticket_id":    true,
		"old_priority": true,
		"created_on":   true,
		"modified_on":  true,
		"created_by":   true,
		"due_date":     true,
		"status":       true,
		"notes":        true,
		"attachments":  true,
	}
)

type ErrorResponse struct {
	StatusCode   int    `json:"status_code"`
	RequestID    string `json:"request_id"`
	ErrorMessage string `json:"error_message"`
	ErrorType    string `json:"error_type"`
}

type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	RequestID  string      `json:"request_id"`
	Data       interface{} `json:"data"`
}

func ValidateStruct(s interface{}) (bool, error) {
	result, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	return result, nil
}

// var validate = validator.New()
// func ValidateStruct[T any](payload T) []*ErrorResponse {
// 	var errors []*ErrorResponse
// 	err := validate.Struct(payload)
// 	if err != nil {
// 		for _, err := range err.(validator.ValidationErrors) {
// 			var element ErrorResponse
// 			element.Field = err.StructNamespace()
// 			element.Tag = err.Tag()
// 			element.Value = err.Param()
// 			errors = append(errors, &element)
// 		}
// 	}
// 	return errors
// }
