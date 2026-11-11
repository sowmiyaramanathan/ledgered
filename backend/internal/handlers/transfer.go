package handlers

import (
	"context"
	"encoding/json"
	"ledgerd/internal/entities/packet"
	"log/slog"
	"net/http"
)

func (h *handler) RecordTransfer(w http.ResponseWriter, r *http.Request) (err error) {
	req := &packet.TransferReq{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Debug("error while parsing", "error", err)
		return err
	}

	err = h.service.RecordTransaction(context.Background(), req)
	if err != nil {
		return
	}

	return
}
