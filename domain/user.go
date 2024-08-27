package domain

import "time"

type User struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required"`

	FirstName     string    `json:"first_name" bson:"first_name"`
	LastName      string    `json:"last_name" bson:"last_name"`
	Verified      bool      `json:"verified" bson:"verified" default:"false"`
	IsAdmin       bool      `json:"is_admin" bson:"is_admin" default:"user"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
	VerifyToken   string    `json:"verify_token,omitempty" bson:"verify_token,omitempty"`
	VerfyTokenExp time.Time `json:"verfy_token_exp,omitempty" bson:"verify_token_exp,omitempty"`
}
