package models

import (
	"context"
	"ledgerd/internal/entities"
	"ledgerd/internal/entities/packet"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type model struct {
	DB *gorm.DB
}

type Model interface {
	GetAccountByID(ctx context.Context, ID uuid.UUID) (account *entities.Account, err error)

	CreateTransaction(ctx context.Context, txn *entities.Transaction) (err error)

	CreateTransferInTxn(ctx context.Context, transfer *packet.TransferData) (err error)
}

func New(db *gorm.DB) Model {
	return &model{
		DB: db,
	}
}
