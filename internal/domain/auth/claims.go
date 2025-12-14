package auth

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/entity"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint        `json:"user_id"`
	Role   entity.Role `json:"role"`
}
