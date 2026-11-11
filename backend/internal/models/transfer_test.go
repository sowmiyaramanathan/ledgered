package models

import (
	"context"
	"ledgerd/apperrors"
	"ledgerd/internal/entities"
	"ledgerd/internal/entities/packet"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("Error loading env file", err)
	}

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Ping())

	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	return db
}

func TestCreateTransferInTxn_Success(t *testing.T) {
	db := setupTestDB(t)
	model := New(db)

	fromAccount := entities.Account{
		ID:        uuid.New(),
		Name:      "Alice",
		Currency:  "INR",
		CreatedAt: time.Now(),
	}

	toAccount := entities.Account{
		ID:        uuid.New(),
		Name:      "Bob",
		Currency:  "INR",
		CreatedAt: time.Now(),
	}

	require.NoError(t, db.Create(&fromAccount).Error)
	require.NoError(t, db.Create(&toAccount).Error)

	txn := &entities.Transaction{
		ID:          uuid.New(),
		Description: "Alice pays Bob",
		CreatedAt:   time.Now(),
	}

	transfer := &packet.TransferData{
		Transaction: txn,
		FromEntry: &entities.Entry{
			ID:        uuid.New(),
			AccountID: fromAccount.ID,
			Amount:    -500,
			CreatedAt: time.Now(),
		},
		ToEntry: &entities.Entry{
			ID:        uuid.New(),
			AccountID: toAccount.ID,
			Amount:    500,
			CreatedAt: time.Now(),
		},
	}

	err := model.CreateTransferInTxn(context.Background(), transfer)
	require.NoError(t, err)

	var txnCount int64
	require.NoError(t,
		db.Model(&entities.Transaction{}).
			Where("id = ?", txn.ID).
			Count(&txnCount).Error,
	)
	require.Equal(t, int64(1), txnCount)

	var entryCount int64
	require.NoError(t,
		db.Model(&entities.Entry{}).
			Where("transaction_id = ?", txn.ID).
			Count(&entryCount).Error,
	)
	require.Equal(t, int64(2), entryCount)
}

func TestCreateTransferInTxn_RollsBackOnSecondEntryFailure(t *testing.T) {
	db := setupTestDB(t)
	model := New(db)
	apperrors.RegisterAppErrors()

	fromAccount := entities.Account{
		ID:        uuid.New(),
		Name:      "Alice",
		Currency:  "INR",
		CreatedAt: time.Now(),
	}

	require.NoError(t, db.Create(&fromAccount).Error)

	// Deliberately do not create this account.
	invalidAccountID := uuid.New()

	txn := &entities.Transaction{
		ID:          uuid.New(),
		Description: "Transfer should fail",
		CreatedAt:   time.Now(),
	}

	transfer := &packet.TransferData{
		Transaction: txn,
		FromEntry: &entities.Entry{
			ID:        uuid.New(),
			AccountID: fromAccount.ID,
			Amount:    -500,
			CreatedAt: time.Now(),
		},
		ToEntry: &entities.Entry{
			ID:        uuid.New(),
			AccountID: invalidAccountID,
			Amount:    500,
			CreatedAt: time.Now(),
		},
	}

	err := model.CreateTransferInTxn(context.Background(), transfer)
	require.Error(t, err)

	// The transaction must have been rolled back.
	var txnCount int64
	require.NoError(t,
		db.Model(&entities.Transaction{}).
			Where("id = ?", txn.ID).
			Count(&txnCount).Error,
	)
	require.Zero(t, txnCount)

	// The first entry must also have been rolled back.
	var entryCount int64
	require.NoError(t,
		db.Model(&entities.Entry{}).
			Where("transaction_id = ?", txn.ID).
			Count(&entryCount).Error,
	)
	require.Zero(t, entryCount)
}
