package auth

import (
	"os"
	"time"
	"errors"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	ErrInvalidToken = errors.New("Invalid or expired token")
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

// Checks token validity and returns associated claims
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
