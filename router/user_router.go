package router

import (
	"example/b/Loan-Tracker-API/controller"
	"example/b/Loan-Tracker-API/infrastructures"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	controller *controller.UserController
	auth       *infrastructures.AuthController
}

func NewUserRouter(uc *controller.UserController) *UserRouter {
	return &UserRouter{
		controller: uc,
		// auth:       ac,
	}
}

func (r *UserRouter) SetupRoutes() {
	router := gin.Default()
	router.POST("/users", r.controller.HandleRegisterUser)
	router.GET("/users/:id/verify", r.controller.HandleVerifyUser)
	router.POST("/users/login", r.controller.HandleAuthenticateUser)

	router.GET("/users/:id", r.auth.ADMINMiddleware(), r.controller.HandleFindUserById)
	router.GET("/users/email/:email", r.auth.ADMINMiddleware(), r.controller.HandleFindUserByEmail)
	router.GET("/users", r.auth.ADMINMiddleware(), r.controller.HandleGetAllUsers)
	router.PUT("/users/:id", r.auth.OWNERMiddleware(), r.controller.HandleUpdateUser)
	router.DELETE("/users/:id", r.auth.OWNERMiddleware(), r.controller.HandleDeleteUser)
	router.PUT("/users/:id/reset-password", r.auth.OWNERMiddleware(), r.controller.HandleResetPassword)
	
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "route not found"})
	})
	router.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "method not found"})
	})
	router.Run(":8080")
}
