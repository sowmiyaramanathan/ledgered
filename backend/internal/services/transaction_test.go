package services

import (
	"context"
	"errors"
	"testing"

	"ledgerd/apperrors"
	"ledgerd/internal/entities"
	"ledgerd/internal/entities/packet"
	"ledgerd/internal/models/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRecordTransaction(t *testing.T) {
	ctx := context.Background()

	fromID := uuid.New()
	toID := uuid.New()

	fromAcc := &entities.Account{
		ID:       fromID,
		Currency: "INR",
	}

	toAcc := &entities.Account{
		ID:       toID,
		Currency: "INR",
	}

	validReq := &packet.TransferReq{
		FromAccountID: fromID,
		ToAccountID:   toID,
		AmountMinor:   10000,
	}

	t.Run("successful transfer", func(t *testing.T) {
		model := mocks.NewMockModel(t)

		model.On("GetAccountByID", ctx, fromID).
			Return(fromAcc, nil).Once()

		model.On("GetAccountByID", ctx, toID).
			Return(toAcc, nil).Once()

		var captured *packet.TransferData

		model.On(
			"CreateTransferInTxn",
			ctx,
			mock.AnythingOfType("*packet.TransferData"),
		).
			Run(func(args mock.Arguments) {
				captured = args.Get(1).(*packet.TransferData)
			}).
			Return(nil).Once()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, validReq)

		require.NoError(t, err)
		require.NotNil(t, captured)

		// Verify transaction
		assert.NotEqual(t, uuid.Nil, captured.Transaction.ID)

		// Verify debit entry
		assert.Equal(t, fromID, captured.FromEntry.AccountID)
		assert.Equal(t, int64(-10000), captured.FromEntry.Amount)

		// Verify credit entry
		assert.Equal(t, toID, captured.ToEntry.AccountID)
		assert.Equal(t, int64(10000), captured.ToEntry.Amount)

		// Verify double-entry invariant
		assert.Equal(
			t,
			int64(0),
			captured.FromEntry.Amount+captured.ToEntry.Amount,
		)

		// Verify unique entry IDs
		assert.NotEqual(
			t,
			captured.FromEntry.ID,
			captured.ToEntry.ID,
		)
	})

	t.Run("nil request", func(t *testing.T) {
		model := mocks.NewMockModel(t)
		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, nil)

		require.Error(t, err)
		model.AssertNotCalled(t, "GetAccountByID", mock.Anything, mock.Anything)
	})

	t.Run("zero amount", func(t *testing.T) {
		model := mocks.NewMockModel(t)
		svc := &service{model: model}

		req := *validReq
		req.AmountMinor = 0

		err := svc.RecordTransaction(ctx, &req)

		require.Error(t, err)
		model.AssertNotCalled(t, "GetAccountByID", mock.Anything, mock.Anything)
	})

	t.Run("negative amount", func(t *testing.T) {
		model := mocks.NewMockModel(t)
		svc := &service{model: model}

		req := *validReq
		req.AmountMinor = -100

		err := svc.RecordTransaction(ctx, &req)

		require.Error(t, err)
		model.AssertNotCalled(t, "GetAccountByID", mock.Anything, mock.Anything)
	})

	t.Run("source account lookup fails", func(t *testing.T) {
		model := mocks.NewMockModel(t)
		dbErr := errors.New("database unavailable")

		model.On("GetAccountByID", ctx, fromID).
			Return(nil, dbErr).Once()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, validReq)

		require.ErrorIs(t, err, dbErr)
		model.AssertNotCalled(t, "CreateTransferInTxn", mock.Anything, mock.Anything)
	})

	t.Run("destination account lookup fails", func(t *testing.T) {
		model := mocks.NewMockModel(t)
		dbErr := errors.New("database unavailable")

		model.On("GetAccountByID", ctx, fromID).
			Return(fromAcc, nil).Once()

		model.On("GetAccountByID", ctx, toID).
			Return(nil, dbErr).Once()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, validReq)

		require.ErrorIs(t, err, dbErr)
		model.AssertNotCalled(t, "CreateTransferInTxn", mock.Anything, mock.Anything)
	})

	t.Run("source account not found", func(t *testing.T) {
		model := mocks.NewMockModel(t)

		model.On("GetAccountByID", ctx, fromID).
			Return(nil, apperrors.DBError).Once()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, validReq)

		// Should return an account-not-found error.
		require.Error(t, err)

		model.AssertNotCalled(t, "CreateTransferInTxn", mock.Anything, mock.Anything)
	})

	t.Run("destination account not found", func(t *testing.T) {
		model := mocks.NewMockModel(t)

		model.On("GetAccountByID", ctx, fromID).
			Return(fromAcc, nil).Once()

		model.On("GetAccountByID", ctx, toID).
			Return(nil, nil).Once()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, validReq)

		require.Error(t, err)

		model.AssertNotCalled(t, "CreateTransferInTxn", mock.Anything, mock.Anything)
	})

	t.Run("same source and destination account", func(t *testing.T) {
		model := mocks.NewMockModel(t)

		sameAccountReq := *validReq
		sameAccountReq.ToAccountID = fromID

		model.On("GetAccountByID", ctx, fromID).
			Return(fromAcc, nil).Twice()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, &sameAccountReq)

		require.Error(t, err)

		model.AssertNotCalled(t, "CreateTransferInTxn", mock.Anything, mock.Anything)
	})

	t.Run("currency mismatch", func(t *testing.T) {
		model := mocks.NewMockModel(t)

		usdAccount := &entities.Account{
			ID:       toID,
			Currency: "USD",
		}

		model.On("GetAccountByID", ctx, fromID).
			Return(fromAcc, nil).Once()

		model.On("GetAccountByID", ctx, toID).
			Return(usdAccount, nil).Once()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, validReq)

		require.Error(t, err)

		model.AssertNotCalled(t, "CreateTransferInTxn", mock.Anything, mock.Anything)
	})

	t.Run("database transaction fails", func(t *testing.T) {
		model := mocks.NewMockModel(t)
		dbErr := errors.New("failed to persist transfer")

		model.On("GetAccountByID", ctx, fromID).
			Return(fromAcc, nil).Once()

		model.On("GetAccountByID", ctx, toID).
			Return(toAcc, nil).Once()

		model.On(
			"CreateTransferInTxn",
			ctx,
			mock.AnythingOfType("*packet.TransferData"),
		).
			Return(dbErr).Once()

		svc := &service{model: model}

		err := svc.RecordTransaction(ctx, validReq)

		require.ErrorIs(t, err, dbErr)
	})
}
