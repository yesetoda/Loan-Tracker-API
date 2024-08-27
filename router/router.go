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
	router.POST("/users/register", r.controller.HandleRegisterUser)
	router.POST("/users/login", r.controller.HandleAuthenticateUser)
	router.GET("/users/verify-email", r.controller.HandleVerifyEmail)

	// TODO:
	router.POST("/users/password-reset/:email", r.controller.HandleRequestResetPassword)
	router.POST("/users/password-update", r.controller.HandleResetPassword)
	router.GET("/users/profile", r.auth.OWNERMiddleware(), r.controller.HandleFindUserById)
	router.GET("/users/token/refresh", r.controller.HandleRefreshToken)

	router.GET("admin/users", r.auth.ADMINMiddleware(), r.controller.HandleGetAllUsers)
	router.GET("admin/users/:id", r.auth.ADMINMiddleware(), r.controller.HandleFindUserById)
	router.DELETE("admin/users/:id", r.auth.ADMINMiddleware(), r.controller.HandleDeleteUser)

	loanRoutes := router.Group("/loans")
	{
		loanRoutes.POST("/", r.auth.AuthenticationMiddleware(), r.controller.ApplyForLoan)
		loanRoutes.GET("/:id", r.auth.AuthenticationMiddleware(), r.controller.ViewLoanStatus)
	}

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(r.auth.ADMINMiddleware())
	{
		adminRoutes.GET("/loans", r.controller.ViewAllLoans)
		adminRoutes.PATCH("/loans/:id/status", r.controller.UpdateLoanStatus)
		adminRoutes.DELETE("/loans/:id", r.controller.DeleteLoan)
	}
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "route not found"})
	})
	router.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "method not found"})
	})
	router.Run(":8080")
}
