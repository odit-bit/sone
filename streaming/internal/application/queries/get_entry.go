package queries

import (
	"context"

	"github.com/odit-bit/sone/streaming/internal/domain"
)

type GetEntry struct {
	Limit  int
	Offset int
}

type GetEntryHandler struct {
	repo domain.EntryRepository
}

func NewGetEntryHandler(repo domain.EntryRepository) *GetEntryHandler {
	return &GetEntryHandler{repo: repo}
}

func (h *GetEntryHandler) GetEntries(ctx context.Context, query GetEntry) (<-chan domain.Entry, error) {
	c := h.repo.List(ctx, query.Offset, query.Limit)
	return c, nil
}
