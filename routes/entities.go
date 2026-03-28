package routes

// import "github.com/google/uuid"

// type UserType string

// const (
// 	Admin        UserType = "admin"
// 	Individual   UserType = "individual"
// 	Corporate    UserType = "corporate"
// 	UnRestricted UserType = ""
// )

// type Role struct {
// 	UserID   uuid.UUID
// 	UserType UserType
// }

// func (r *Role) IsAdmin() bool {
// 	return r.UserType == Admin
// }

// func (r *Role) IsCorporate() bool {
// 	return r.UserType == Corporate
// }

// func (r *Role) IsIndividual() bool {
// 	return r.UserType == Individual
// }

// func (r *Role) IsUser() bool {
// 	return r.UserType == Individual || r.UserType == Corporate
// }

// func (r *Role) IsUnRestricted() bool {
// 	return r.IsAdmin() || r.IsIndividual() || r.IsCorporate()
// }

// func (r *Role) IsUnknownRole() bool {
// 	return !r.IsUnRestricted()
// }
