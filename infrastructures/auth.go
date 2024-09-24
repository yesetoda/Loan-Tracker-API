package infrastructures

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/yesetoda/b/Loan-Tracker-API/config"
	"github.com/yesetoda/b/Loan-Tracker-API/domain"
	"github.com/yesetoda/b/Loan-Tracker-API/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthController struct {
	userRepo repository.GeneralRepository
}

func NewAuthController(userRepo repository.GeneralRepository) *AuthController {
	return &AuthController{
		userRepo: userRepo,
	}
}

func (ac *AuthController) AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "unexpected error"})
				c.Abort()
			}
		}()
		claims, err := GetClaims(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func (ac *AuthController) ADMINMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "unexpected error"})
				c.Abort()
			}
		}()
		claims, err := GetClaims(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if claims.IsAdmin {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, errors.New("invalid token"))
	}
}

func (ac *AuthController) USERMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "unexpected error"})
				c.Abort()
			}
		}()
		claim, err := GetClaims(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if claim.IsActive {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, errors.New("invalid token"))
		c.Abort()
	}

}
func (ac *AuthController) OWNERMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "unexpected error"})
				c.Abort()
			}
		}()
		fmt.Println("this is the owner middleware")
		claim, err := GetClaims(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		username := c.Param("username")
		fmt.Println("this is the claim,username", claim, username)
		fmt.Println("start calling find user by username ")
		user, err := ac.userRepo.FindUserBy(username)
		fmt.Println("end calling find user by username ")

		fmt.Println("this is the username", username, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		fmt.Println("this is the user", user, username)
		if user.Username == claim.Username {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, errors.New("unauthorized,neither an admin nor an author"))
	}
}
func GetClaims(c *gin.Context) (*domain.Claims, error) {
	config_domain, err := config.LoadConfig()
	if err != nil {
		return &domain.Claims{}, err
	}

	var jwtSecret = []byte(config_domain.Jwt.JwtKey)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return &domain.Claims{}, errors.New("missing authorization header")
	}

	TokenString := strings.Split(authHeader, " ")
	if len(TokenString) != 2 || TokenString[0] != "Bearer" {
		return &domain.Claims{}, errors.New("invalid token format")
	}
	tokenString := TokenString[1]

	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return &domain.Claims{}, err
	}
	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, err
	}
	return &domain.Claims{}, errors.New("invalid token")
}
