package models

import (
	"github.com/google/uuid"
	"time"
)

type UserRegistration struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	BirthDate time.Time `json:"birthDate"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
}
