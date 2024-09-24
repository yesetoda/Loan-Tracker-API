package repository

import (
	"github.com/yesetoda/b/Loan-Tracker-API/domain"
)

type GeneralRepository interface {
	CreateUser(domain.User) (domain.User, error)
	FindUserByEmail(string) (domain.User, error)
	FindUserBy(id string) (domain.User, error)
	UpdateUser(id string, user domain.User) error
	DeleteUser(id string) (string error)
	VerifyUser(id string) (string error)
	AuthenticateUser(id, password string) (string error)
	ListAllUsers() ([]domain.User,error)
	ResetPassword(id, password string) (string error)

	CreateLoan(loan domain.Loan) (domain.Loan, error)
	FindLoanByID(id string) (domain.Loan, error)
	FindLoansByUserID(userID string) ([]domain.Loan, error)
	FindAllLoans(status string, order string) ([]domain.Loan, error)
	UpdateLoanStatus(id string, status string) (string, error)
	DeleteLoan(id string) (string, error)
}
