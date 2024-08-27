package controller

import (
	"example/b/Loan-Tracker-API/domain"
	"example/b/Loan-Tracker-API/usecase"

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

func (c *UserController) HandleFindUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.usecase.FindUserById(id)
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
	err := c.usecase.DeleteUser(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (c *UserController) HandleVerifyUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.usecase.VerifyUser(id)
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
	err := c.usecase.AuthenticateUser(user.ID, user.Password)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User authenticated successfully"})
}

func (c *UserController) HandleResetPassword(ctx *gin.Context) {
	id := ctx.Param("id")
	var password domain.Password
	if err := ctx.ShouldBindJSON(&password); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := c.usecase.ResetPassword(id, password.NewPassword)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Password reset successfully"})
}
