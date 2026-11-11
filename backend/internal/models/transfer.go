package models

import (
	"context"
	"ledgerd/apperrors"
	"ledgerd/internal/entities/packet"
	"log/slog"
)

func (m *model) CreateTransferInTxn(ctx context.Context, transfer *packet.TransferData) (err error) {
	tx := m.DB.WithContext(ctx).Begin()

	if tx.Error != nil {
		slog.Error("error while creating a transaction")
		return apperrors.DBError.SetError(tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			slog.Error("panic occurred in CreateTransferInTxn", "panic", r)
			err = apperrors.DBError.SetDetailParameters("panic in CreateTransferInTxn")
		}
	}()

	err = tx.Create(&transfer.Transaction).Error
	if err != nil {
		tx.Rollback()
		slog.Warn("error while creating transaction", "error", err, "transaction", transfer.Transaction)
		return apperrors.DBError.SetError(err)
	}

	// Assign the transaction ID to the entries before inserting them
	transfer.FromEntry.TransactionID = transfer.Transaction.ID
	transfer.ToEntry.TransactionID = transfer.Transaction.ID

	err = tx.Create(&transfer.FromEntry).Error
	if err != nil {
		tx.Rollback()
		slog.Warn("error while creating entry for from account", "error", err, "fromAccountID", transfer.FromEntry.AccountID)
		return apperrors.DBError.SetError(err)
	}

	err = tx.Create(&transfer.ToEntry).Error
	if err != nil {
		tx.Rollback()
		slog.Warn("error while creating entry for to account", "error", err, "toAccountID", transfer.ToEntry.AccountID)
		return apperrors.DBError.SetError(err)
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		slog.Error("transaction commit failed in CreateTransferInTxn", "error", commitErr)
		return apperrors.DBError.SetError(commitErr)
	}

	return
}
