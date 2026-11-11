package entities

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(120);not null"`
	Currency  string    `json:"currency" gorm:"type:char(3);not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
}
