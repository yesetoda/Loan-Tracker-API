package repository

import "example/b/Loan-Tracker-API/domain"

type GeneralRepository interface {
	CreateUser(domain.User) (domain.User, error)
	FindUserByEmail(string) (domain.User, error)
	FinduserById(id string) (domain.User, error)
	UpdateUser(id string, user domain.User) error
	DeleteUser(id string) error
	VerifiyUser(id string) error
	AuthenticateUser(id, password string) error
	ListAllUsers() []domain.User
	ResetPassword(id, password string) error
}
