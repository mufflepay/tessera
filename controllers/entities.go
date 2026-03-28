package controllers

import "time"

type UserType string

const (
	Admin      UserType = "admin"
	Individual UserType = "individual"
	Corporate  UserType = "corporate"
)

type UserInfo struct {
	UserID    string
	Role      UserType
	UserAgent []byte
	IPAddress string
	ExpiresOn time.Time
}

func (u *UserInfo) IsAdmin() bool {
	return u.Role == Admin
}

func (u *UserInfo) IsNotAdmin() bool {
	return !u.IsAdmin()
}

func (u *UserInfo) IsCorporate() bool {
	return u.Role == Corporate
}

func (u *UserInfo) IsIndividual() bool {
	return u.Role == Individual
}

func (u *UserInfo) IsAdminOrUser() bool {
	return u.IsAdmin() || u.IsIndividual() || u.IsCorporate()
}

func (u *UserInfo) IsNotAdminOrCorporate() bool {
	return !u.IsAdmin() || !u.IsCorporate()
}

func (u *UserInfo) IsNotAdminOrUser() bool {
	return !u.IsAdminOrUser()
}
