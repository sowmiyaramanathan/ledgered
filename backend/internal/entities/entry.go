package entities

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey"`
	TransactionID uuid.UUID   `gorm:"type:uuid;not null;index"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID;references:ID"`
	AccountID     uuid.UUID   `gorm:"type:uuid;not null;index"`
	Account       Account     `gorm:"foreignKey:AccountID;references:ID"`
	Amount        int64       `gorm:"type:bigint;not null"`
	CreatedAt     time.Time   `json:"createdAt" gorm:"not null"`
}
