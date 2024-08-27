package domain

import "github.com/golang-jwt/jwt"

type Claims struct {
	ID       string `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	jwt.StandardClaims
}
