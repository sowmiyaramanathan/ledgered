package services

import (
	"context"
	"ledgerd/apperrors"
	"ledgerd/internal/entities"
	"ledgerd/internal/entities/packet"
	"math"
	"time"

	"github.com/google/uuid"
)

func (s *service) RecordTransaction(ctx context.Context, req *packet.TransferReq) (err error) {

	if req == nil || req.AmountMinor <= 0 {
		return apperrors.SInvalidAmountErr
	}

	fromAcc, err := s.model.GetAccountByID(ctx, req.FromAccountID)
	if err != nil {
		return
	}

	toAcc, err := s.model.GetAccountByID(ctx, req.ToAccountID)
	if err != nil {
		return
	}

	if fromAcc == nil || toAcc == nil {
		return apperrors.SAccountNotFoundErr
	}

	if fromAcc.ID == toAcc.ID {
		return apperrors.SSameAccoutErr
	}

	if fromAcc.Currency != toAcc.Currency {
		return apperrors.SCurrencyMismatchErr
	}

	if req.AmountMinor < 0 {
		return apperrors.SEntryBalanceErr
	}

	transferData := &packet.TransferData{
		FromEntry: &entities.Entry{
			ID:        uuid.New(),
			AccountID: fromAcc.ID,
			Amount:    -req.AmountMinor,
			CreatedAt: time.Now(),
		},
		ToEntry: &entities.Entry{
			ID:        uuid.New(),
			AccountID: toAcc.ID,
			Amount:    req.AmountMinor,
			CreatedAt: time.Now(),
		},
		Transaction: &entities.Transaction{
			ID: uuid.New(),
		},
	}
	if err := validateTransfer(transferData); err != nil {
		return err
	}

	err = s.model.CreateTransferInTxn(ctx, transferData)
	if err != nil {
		return
	}

	return
}

func validateTransfer(data *packet.TransferData) error {
	if data == nil || data.FromEntry == nil ||
		data.ToEntry == nil || data.Transaction == nil {
		return apperrors.SEntryBalanceErr
	}

	from := data.FromEntry
	to := data.ToEntry

	if from.AccountID == to.AccountID {
		return apperrors.SSameAccoutErr
	}

	if from.Amount == math.MinInt64 || -from.Amount != to.Amount {
		return apperrors.SEntryBalanceErr
	}

	if from.Amount+to.Amount != 0 {
		return apperrors.SEntryBalanceErr
	}

	return nil
}
