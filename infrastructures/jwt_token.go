package infrastructures

import (
	"errors"
	"example/b/Loan-Tracker-API/domain"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *domain.User, password, jwtSecret string) (string, string, error) {
	// Compare provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid username or password")
	}

	// Generate Access Token
	accessToken, err := createJWTToken(user, jwtSecret, 30*time.Minute)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken, err := createJWTToken(user, jwtSecret, 3*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func createJWTToken(user *domain.User, jwtSecret string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &domain.Claims{
		ID:       user.ID,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
		IsActive: user.Verified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
