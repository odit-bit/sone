package rpc

import (
	"context"

	"github.com/odit-bit/sone/ingress/internal/domain"
	"github.com/odit-bit/sone/streaming/streamingpb"
)

type Streams struct {
	ls streamingpb.LiveStreamClient
}

func NewLiveStreams(cli streamingpb.LiveStreamClient) *Streams {
	return &Streams{
		ls: cli,
	}
}

func (s *Streams) StartStream(ctx context.Context, key string) (*domain.StreamInfo, error) {
	resp, err := s.ls.Start(ctx, &streamingpb.StartRequest{Key: key})
	if err != nil {
		return nil, err
	}
	return &domain.StreamInfo{ID: resp.GetId(), Title: resp.Title}, nil
}

func (s *Streams) EndStream(ctx context.Context, key string) (*domain.StreamInfo, error) {
	info, err := s.ls.End(ctx, &streamingpb.EndRequest{Key: key})
	if err != nil {
		return nil, err
	}

	return &domain.StreamInfo{ID: info.Id, Title: info.Title}, nil
}
