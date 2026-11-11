package packet

import (
	"ledgerd/internal/entities"

	"github.com/google/uuid"
)

type TransferReq struct {
	FromAccountID uuid.UUID
	ToAccountID   uuid.UUID
	AmountMinor   int64
}

type TransferData struct {
	FromEntry   *entities.Entry
	ToEntry     *entities.Entry
	Transaction *entities.Transaction
}
