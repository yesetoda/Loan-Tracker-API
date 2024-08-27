package usecase

import (
	"crypto/rand"
	"example/b/Loan-Tracker-API/domain"
	"example/b/Loan-Tracker-API/infrastructures"
	"example/b/Loan-Tracker-API/infrastructures/password_service"

	"example/b/Loan-Tracker-API/repository"
	"math/big"
)

type GeneralUsecase struct {
	UserRepo repository.GeneralRepository
}

func NewUsecase(ur repository.GeneralRepository) GeneralUsecase {
	return GeneralUsecase{
		UserRepo: ur,
	}
}

func (uc *GeneralUsecase) RegisterUser(user domain.User) error {
	user.Password, _ = password_service.HashPassword(user.Password)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	confirmationToken := make([]byte, 64)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < 64; i++ {
		num, _ := rand.Int(rand.Reader, charsetLength)
		confirmationToken[i] = charset[num.Int64()]
	}
	user.VerifyToken = string(confirmationToken)
	user, err := uc.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		link := "localhost:8080/users/accountVerification/?email=" + user.Email + "&token=" + string(confirmationToken)
		err = infrastructures.SendEmail(user.Email, "Registration Confirmation", "This sign up Confirmation email to verify: ", link)
		if err != nil {
			return err
		}
	} else {
		infrastructures.SendEmail(user.Email, "Welcome! You Are Our First Admin", "Congratulations! As the first user to join our site, you have been automatically granted admin privileges. Thank you for being an early supporter.", "")
	}
	if err != nil {
		return err
	}
	return err
}
func (uc *GeneralUsecase) ListAllUsers() []domain.User {
	return uc.UserRepo.ListAllUsers()
}
func (uc *GeneralUsecase) FindUserById(id string) (domain.User, error) {
	return uc.UserRepo.FinduserById(id)
}

func (uc *GeneralUsecase) FinduserByEmail(email string) (domain.User, error) {
	return uc.UserRepo.FindUserByEmail(email)
}

func (uc *GeneralUsecase) UpdateUser(id string, userData domain.User) error {
	return uc.UserRepo.UpdateUser(id, userData)
}

func (uc *GeneralUsecase) DeleteUser(id string) error {
	return uc.UserRepo.DeleteUser(id)
}

func (uc *GeneralUsecase) VerifyUser(id string) error {
	return uc.UserRepo.VerifiyUser(id)
}
func (uc *GeneralUsecase) AuthenticateUser(id, password string) error {
	return uc.UserRepo.AuthenticateUser(id, password)
}

func (uc *GeneralUsecase) ResetPassword(id, password string) error {
	return uc.UserRepo.ResetPassword(id, password)
}
