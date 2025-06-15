package auth

import (
	"os"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

// claims represents token encoded data
type Claims struct {
	UserID uuid.UUID `json:user_id`
	jwt.RegisteredClaims
}

// Returns a signed JWT for a given user
func GenerateToken(user *models.User) (string, error) {
	claims := Claims {
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims {
			ExpiresAt: 	jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:	jwt.NewNumericDate(time.Now()),	
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
