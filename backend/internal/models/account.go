package models

import (
	"context"
	"ledgerd/apperrors"
	"ledgerd/internal/entities"
	"log/slog"

	"github.com/google/uuid"
)

func (m *model) GetAccountByID(ctx context.Context, ID uuid.UUID) (account *entities.Account, err error) {
	account = &entities.Account{}

	err = m.DB.WithContext(ctx).Model(&entities.Account{}).Where("id = ?", ID).First(&account).Error
	if err != nil {
		slog.Warn("error fetching account from db", "error", err)
		return account, apperrors.DBError.SetDetailParameters("error fetching account from DB")
	}

	return
}
