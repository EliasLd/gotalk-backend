package auth

import (
	"errors"
	"os"
	"time"

	"github.com/EliasLd/gotalk-backend/internal/models"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var env_file_path

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)
