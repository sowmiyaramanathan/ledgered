package handlers

import (
	"ledgerd/internal/services"
	"net/http"
)

type handler struct {
	service services.Service
}

type Handler interface {
	RecordTransfer(w http.ResponseWriter, r *http.Request) (err error)
}

func New(s *services.Service) Handler {
	return &handler{
		service: *s,
	}
}
