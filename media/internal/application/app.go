package application

import (
	"context"
	"io"

	"github.com/odit-bit/sone/media/internal/domain"
)

type Media struct {
	segments domain.SegmentRepository
}

func New(segments domain.SegmentRepository) *Media {
	return &Media{segments: segments}
}

type GetSegmentArgs struct {
	Name string
}

func (app *Media) Get(ctx context.Context, arg GetSegmentArgs) (*domain.Segment, error) {
	segment, err := app.segments.GetVideoSegment(ctx, arg.Name)
	if err != nil {
		return nil, err
	}
	return segment, nil
}

type PutSegmentArgs struct {
	Key  string
	Body io.Reader
}

func (app *Media) Put(ctx context.Context, arg PutSegmentArgs) error {
	return app.segments.InsertVideoSegment(ctx, arg.Key, arg.Body)
}
