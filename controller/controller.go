package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yesetoda/b/Loan-Tracker-API/config"
	"github.com/yesetoda/b/Loan-Tracker-API/domain"
	"github.com/yesetoda/b/Loan-Tracker-API/infrastructures"
	"github.com/yesetoda/b/Loan-Tracker-API/infrastructures/password_service"
	"github.com/yesetoda/b/Loan-Tracker-API/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase *usecase.GeneralUsecase
}

func NewUserController(usecase *usecase.GeneralUsecase) *UserController {
	return &UserController{
		usecase: usecase,
	}
}

func (c *UserController) HandleRegisterUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := c.usecase.RegisterUser(user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User registered successfully"})
}

func (c *UserController) HandleFindUserByUsername(ctx *gin.Context) {
	fmt.Println("this is the Handle find user by username")

	username := ctx.Param("username")
	fmt.Println("this is the username", username)
	user, err := c.usecase.FindUserByUsername(username)
	fmt.Println("this is the user and error", user, err)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, user)
}
func (c *UserController) HandleFindUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.usecase.FinduserByEmail(email)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, user)
}

func (c *UserController) HandleGetAllUsers(ctx *gin.Context) {
	users := c.usecase.ListAllUsers()
	ctx.JSON(200, users)
}

func (c *UserController) HandleRefreshToken(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	the_user, err := c.usecase.FinduserByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no such user"})
		return
	}
	if password_service.CheckPasswordHash(user.Password, the_user.Password) != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load config"})
		return
	}
	newAT, newRT, err := infrastructures.GenerateToken(&the_user, user.Password, cfg.Jwt.JwtKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	ctx.JSON(200, gin.H{"access_token": newAT, "refresh_token": newRT})
}

func (c *UserController) HandleUpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := c.usecase.UpdateUser(id, user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, user)
}

func (c *UserController) HandleDeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	claims, err := infrastructures.GetClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}
	if claims.ID == id && claims.IsAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "can't delete admin"})
		return
	}
	err = c.usecase.DeleteUser(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (c *UserController) HandleVerifyEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	token := ctx.Query("token")
	fmt.Println("this is the email address", email, token)
	err := c.usecase.VerifyUser(email, token)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User verified successfully"})
}

func (c *UserController) HandleAuthenticateUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("this is the user", user)
	at, rt, err := c.usecase.AuthenticateUser(user.Email, user.Password)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"access_token": at, "refresh_token": rt})
}

func (c *UserController) HandleRequestResetPassword(ctx *gin.Context) {
	email := ctx.Param("email")
	messagge, err := c.usecase.PasswordReset(email)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error(), "message": messagge})
		return
	}
	ctx.JSON(200, gin.H{"message": messagge})
}

func (c *UserController) HandleResetPassword(ctx *gin.Context) {
	email := ctx.Query("email")
	token := ctx.Query("token")
	newPassword := ctx.Query("password")
	messagge, err := c.usecase.PasswordUpdate(email, token, newPassword)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error(), "message": messagge})
		return
	}
	ctx.JSON(200, gin.H{"message": messagge})
}

func (c *UserController) ApplyForLoan(ctx *gin.Context) {
	var loan domain.Loan
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims, err := infrastructures.GetClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	loan.UserID = claims.ID

	loan, err = c.usecase.ApplyForLoan(loan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply for loan"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"created loanID": loan.ID})
}

// ViewLoanStatus retrieves the status of a specific loan
func (c *UserController) ViewLoanStatus(ctx *gin.Context) {
	loanID := ctx.Param("id")
	fmt.Println("this is the id", loanID)
	loan, err := c.usecase.ViewLoanStatus(loanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loan status"})
		return
	}

	ctx.JSON(http.StatusOK, loan)
}

// ViewAllLoans retrieves all loan applications (Admin only)
func (c *UserController) ViewAllLoans(ctx *gin.Context) {
	status := ctx.Query("status")
	order := ctx.Query("order")

	loans, err := c.usecase.ViewAllLoans(status, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loans"})
		return
	}

	ctx.JSON(http.StatusOK, loans)
}

// UpdateLoanStatus updates the status of a loan application (Admin only)
func (c *UserController) UpdateLoanStatus(ctx *gin.Context) {
	loanID := ctx.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := c.usecase.UpdateLoanStatus(loanID, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan status", "message": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan status updated successfully"})
}

// DeleteLoan deletes a specific loan application (Admin only)
func (c *UserController) DeleteLoan(ctx *gin.Context) {
	loanID := ctx.Param("id")
	message, err := c.usecase.DeleteLoan(loanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete loan", "message": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}
