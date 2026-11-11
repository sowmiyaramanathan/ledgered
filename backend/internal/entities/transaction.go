package entities

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Description string    `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"not null"`
}
