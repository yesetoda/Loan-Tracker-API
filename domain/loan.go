package domain

import (
	"time"
)

type Loan struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Amount    float64   `json:"amount" bson:"amount" binding:"required"`
	Term      int       `json:"term" bson:"term"`
	Status    string    `json:"status" bson:"status"` // Status of the loan: "pending", "approved", "rejected"
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
