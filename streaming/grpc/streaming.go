package grpc

import (
	"context"
	"fmt"

	"github.com/odit-bit/sone/streaming/internal/app"
	"github.com/odit-bit/sone/streaming/streamingpb"
)

var _ streamingpb.LiveStreamServer = (*Server)(nil)

type Server struct {
	app *app.Application
	streamingpb.UnimplementedLiveStreamServer
}

func NewServer(app *app.Application) *Server {
	return &Server{
		app:                           app,
		UnimplementedLiveStreamServer: streamingpb.UnimplementedLiveStreamServer{},
	}
}

// Insert implements streamingpb.LiveStreamServer.
func (s *Server) Insert(ctx context.Context, req *streamingpb.InsertRequest) (*streamingpb.LiveStreamInfo, error) {
	res, err := s.app.InsertStream(ctx, req.Token, app.InsertStreamOpt{Title: req.Title})
	if err != nil {
		return nil, err
	}
	return &streamingpb.LiveStreamInfo{Id: res.ID, Title: res.Title, IsLive: int32(res.IsLive), Key: res.Key}, nil
}

func (s *Server) Start(ctx context.Context, req *streamingpb.StartRequest) (*streamingpb.StartResponse, error) {
	ls, err := s.app.StartStream(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &streamingpb.StartResponse{Id: ls.ID, Title: ls.Title, IsLive: int32(ls.IsLive)}, nil
}

func (s *Server) End(ctx context.Context, req *streamingpb.EndRequest) (*streamingpb.LiveStreamInfo, error) {
	res, err := s.app.EndStream(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &streamingpb.LiveStreamInfo{Id: res.ID, Title: res.Title, IsLive: int32(res.IsLive)}, nil
}

func (s *Server) List(r *streamingpb.ListRequest, w streamingpb.LiveStream_ListServer) error {
	c, err := s.app.ListStream(w.Context(), r.Id)
	if err != nil {
		return err
	}

	for {
		ls, ok := c.Next()
		if !ok {
			break
		}
		fmt.Println(ls)
		if err := w.Send(&streamingpb.LiveStreamInfo{
			Id:     ls.ID,
			Title:  ls.Title,
			IsLive: int32(ls.IsLive),
		}); err != nil {
			return err
		}
	}

	return c.Err()
}
