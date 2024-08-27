package usecase

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"example/b/Loan-Tracker-API/config"
	"example/b/Loan-Tracker-API/domain"
	"example/b/Loan-Tracker-API/infrastructures"
	"example/b/Loan-Tracker-API/infrastructures/password_service"
	"example/b/Loan-Tracker-API/repository"
)

type GeneralUsecase struct {
	Repo repository.GeneralRepository
}

func NewUsecase(ur repository.GeneralRepository) GeneralUsecase {
	return GeneralUsecase{
		Repo: ur,
	}
}

func generateConfirmationToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	confirmationToken := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range confirmationToken {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		confirmationToken[i] = charset[num.Int64()]
	}

	return string(confirmationToken), nil
}
func (uc *GeneralUsecase) RegisterUser(user domain.User) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return err
	}
	hashedPassword, err := password_service.HashPassword(user.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}
	user.Password = hashedPassword

	confirmationToken, err := generateConfirmationToken(64)
	if err != nil {
		log.Printf("Failed to generate confirmation token: %v", err)
		return err
	}
	user.VerifyToken = confirmationToken

	user, err = uc.Repo.CreateUser(user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return err
	}

	// Construct the verification link
	link := fmt.Sprintf("%s/users/verify-email/?email=%s&token=%s", cfg.Port, user.Email, confirmationToken)

	// Send the verification email
	var emailSubject, emailBody string
	if user.IsAdmin {
		emailSubject = "Welcome! You Are Our First Admin"
		emailBody = "Congratulations! As the first user to join our site, you have been automatically granted admin privileges. Thank you for being an early supporter."
	} else {
		emailSubject = "Registration Confirmation"
		emailBody = "This is a sign-up confirmation email to verify your account."
	}

	if err = infrastructures.SendEmail(user.Email, emailSubject, emailBody, link); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func (uc *GeneralUsecase) ListAllUsers() []domain.User {
	users, _ := uc.Repo.ListAllUsers()
	return users
}
func (uc *GeneralUsecase) FindUserById(id string) (domain.User, error) {
	return uc.Repo.FindUserById(id)
}

func (uc *GeneralUsecase) FinduserByEmail(email string) (domain.User, error) {
	return uc.Repo.FindUserByEmail(email)
}

func (uc *GeneralUsecase) UpdateUser(id string, userData domain.User) error {
	return uc.Repo.UpdateUser(id, userData)
}

func (uc *GeneralUsecase) DeleteUser(id string) error {
	return uc.Repo.DeleteUser(id)
}

func (uc *GeneralUsecase) VerifyUser(email, token string) error {
	// Find the user by email
	user, err := uc.Repo.FindUserByEmail(email)
	if err != nil {
		log.Printf("Failed to find user with email %s: %v", email, err)
		return err
	}
	fmt.Println("this is the user", user, user.Verified)
	// Check if the user is already verified
	if user.Verified {
		log.Printf("User with email %s is already verified", email)
		return nil
	}

	// Validate the token
	if user.VerifyToken != token {
		log.Printf("Invalid token for user with email %s", email)
		return errors.New("invalid token")
	}

	// Mark the user as verified
	user.Verified = true
	user.VerifyToken = "" // Clear the token once it's used

	// Update the user in the database
	err = uc.Repo.UpdateUser(user.ID, user)
	if err != nil {
		log.Printf("Failed to verify user with email %s: %v", email, err)
		return err
	}

	log.Printf("User with email %s successfully verified", email)
	return nil
}

func (uc *GeneralUsecase) AuthenticateUser(email, password string) (string, string, error) {
	fmt.Println("this is get by email", email, password)
	user, err := uc.Repo.FindUserByEmail(email)
	fmt.Println(user, err)
	if err != nil {
		log.Printf("Failed to find user with email %s: %v", email, err)
		return "", "", err
	}

	if password_service.CheckPasswordHash(password, user.Password) != nil {
		log.Printf("Authentication failed for user with email %s: incorrect password", email)
		return "", "", errors.New("incorrect password")
	}
	fmt.Println("password hash is correct")
	config_domain, err := config.LoadConfig()
	if err != nil {
		return "", "", err
	}
	fmt.Println(config_domain.Jwt.JwtKey)
	accesstoken, refreshtokeng, err := infrastructures.GenerateToken(&user, password, config_domain.Jwt.JwtKey)
	fmt.Println("this is the tokens", accesstoken, refreshtokeng, err)
	if err != nil {
		return "", "", err
	}
	log.Printf("User with email %s successfully authenticated", email)
	return accesstoken, refreshtokeng, err
}

func (uc *GeneralUsecase) PasswordReset(email string) (string, error) {

	user, err := uc.FinduserByEmail(email)
	if err != nil {
		return "", err
	}
	if !user.Verified {
		return "", errors.New("account not activated")
	}
	confirmationToken, err := generateConfirmationToken(64)
	if err != nil {
		return "", err
	}
	config_domain, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	link := config_domain.Domain + "/users/resetPassword/?email=" + user.Email + "&token=" + string(confirmationToken)
	err = infrastructures.SendEmail(user.Email, "Password Reset", "This is the password reset link: ", link)
	if err != nil {
		return "", err
	}

	return "Password reset token sent to your email", nil
}

func (uc *GeneralUsecase) PasswordUpdate(email, token, password string) (string, error) {

	user, err := uc.FinduserByEmail(email)
	if err != nil {
		return "", err
	}
	if !user.Verified {
		return "", errors.New("account not activated")
	}
	if user.VerifyToken == token {
		if user.VerfyTokenExp.Before(time.Now()) {
			return "Token has expired", errors.New("token expired")
		}
		user.Password, _ = password_service.HashPassword(password)
		err := uc.Repo.UpdateUser(user.ID, user)
		if err != nil {
			return "password has not been updated", err
		}
		return "Password reset successful", nil
	}
	return "Invalid token", errors.New("invalid token")
}

func (uc *GeneralUsecase) ApplyForLoan(loan domain.Loan) (domain.Loan, error) {
	loan.Status = "pending"
	loan, err := uc.Repo.CreateLoan(loan)
	if err != nil {
		return loan, err
	}
	return loan, nil
}

func (uc *GeneralUsecase) ViewLoanStatus(loanID string) (*domain.Loan, error) {
	loan, err := uc.Repo.FindLoanByID(loanID)
	fmt.Println("this si the loan status", loan, err)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (uc *GeneralUsecase) ViewAllLoans(status, order string) ([]domain.Loan, error) {
	loans, err := uc.Repo.FindAllLoans(status, order)
	if err != nil {
		log.Printf("Failed to retrieve all loans: %v", err)
		return nil, err
	}
	return loans, nil
}

func (uc *GeneralUsecase) UpdateLoanStatus(loanID string, status string) (string, error) {
	message, err := uc.Repo.UpdateLoanStatus(loanID, status)
	if err != nil {
		return message, err
	}
	return message, nil
}

func (uc *GeneralUsecase) DeleteLoan(loanID string) (string, error) {
	message, err := uc.Repo.DeleteLoan(loanID)
	if err != nil {
		return message, err
	}
	return message, nil
}
