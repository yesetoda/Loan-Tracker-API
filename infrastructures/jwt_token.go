package infrastructures

import (
	"errors"
	"example/b/Loan-Tracker-API/domain"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *domain.User, pwd string) (string, string, error) {
	jwtSecret := []byte("secret")

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)) != nil {
		return "", "", errors.New("invalid username or password")
	}

	expirationTime := time.Now().Add(30 * time.Minute)
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

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	expirationTime = time.Now().Add(3 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
