package commands

import (
	"context"

	"github.com/odit-bit/sone/streaming/internal/domain"
)

type CreateStreamKeyCMD struct {
	Name string
}

type CreateStreamKeyHandler struct {
	repo domain.StreamKeyRepository
}

func NewCreateStreamKeyHandler(repo domain.StreamKeyRepository) *CreateStreamKeyHandler {
	return &CreateStreamKeyHandler{repo: repo}
}

func (h *CreateStreamKeyHandler) CreateStreamKey(ctx context.Context, cmd CreateStreamKeyCMD) (domain.Key, error) {
	key := domain.CreateKey()
	err := h.repo.Save(ctx, key)
	return key, err
}
