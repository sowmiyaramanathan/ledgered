package services

import (
	"context"
	"ledgerd/internal/entities/packet"
	models "ledgerd/internal/models"
)

type service struct {
	model models.Model
}

type Service interface {
	RecordTransaction(ctx context.Context, req *packet.TransferReq) (err error)
}

func New(m *models.Model) Service {
	return &service{
		model: *m,
	}
}
