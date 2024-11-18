package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/odit-bit/sone/streaming/internal/domain"
)

type GetSegment struct {
	Name string
}

type GetSegmentHandler struct {
	repo domain.SegmentRepository
}

func NewGetSegmentHandler(repo domain.SegmentRepository) *GetSegmentHandler {
	return &GetSegmentHandler{
		repo: repo,
	}
}

func (h *GetSegmentHandler) GetSegment(ctx context.Context, query GetSegment) (*domain.Segment, error) {
	seg, err := h.repo.Get(ctx, query.Name)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("get segment query"))
	}
	return seg, nil
}
