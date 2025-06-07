package models

import (
	"time"
	
	"github.com/google/uuid"
)

type User struct {
	ID		uuid.UUID 	`db:"id"`
	Username	string		`db:"username"`
	Password	string		`db:"password_hash"`
	CreatedAt	time.Time	`db:"created_at"`
	UpdatedAt	time.Time	`db:"updated_at"`	
}
