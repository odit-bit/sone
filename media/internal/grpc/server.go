package grpc

import (
	"io"

	"github.com/odit-bit/sone/media/internal/application"
	"github.com/odit-bit/sone/media/mediapb"
	"github.com/valyala/bytebufferpool"
	"google.golang.org/grpc"
)

type server struct {
	app *application.Media
	mediapb.UnimplementedMediaServiceServer
}

var _ mediapb.MediaServiceServer = (*server)(nil)

func RegisterServer(app *application.Media, registrar grpc.ServiceRegistrar) error {
	srv := &server{app: app}
	mediapb.RegisterMediaServiceServer(registrar, srv)
	return nil
}

func (s *server) GetSegment(req *mediapb.SegmentRequest, res mediapb.MediaService_GetSegmentServer) error {
	seg, err := s.app.Get(res.Context(), application.GetSegmentArgs{Name: req.Path})
	if err != nil {
		return err
	}
	defer seg.Body.Close()

	// candidate for optimized ??
	// buf := make([]byte, 1024*1024)
	p := bytebufferpool.Get()
	defer bytebufferpool.Put(p)
	buf := p.B

	for {
		n, err := seg.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			// finish read
			break
		}

		if err := res.Send(&mediapb.SegmentResponse{
			BodyChunk: buf[:n],
		}); err != nil {
			return err
		}

	}
	return nil
}

// func (s *server) PutSegment(ctx context.Context, req *mediapb.PutSegmentRequest) (*mediapb.PutSegmentResponse, error) {
// 	buf := buffer.Pool.Get()
// 	_, err := buf.Write(req.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := s.app.Put(ctx, application.PutSegmentArgs{
// 		Key:  req.Path,
// 		Body: buf,
// 	}); err != nil {
// 		return nil, err
// 	}
// 	buf.Reset()
// 	buffer.Pool.Put(buf)
// 	return &mediapb.PutSegmentResponse{}, nil
// }

// // ListQuery implements mediapb.MediaServiceServer.
// func (s *server) ListQuery(req *mediapb.ListMediaRequest, stream mediapb.MediaService_ListQueryServer) error {
// 	ctx, cancel := context.WithCancel(stream.Context())
// 	defer cancel()

// 	entry, err := s.app.GetEntries(ctx, queries.GetEntry{
// 		Limit:  int(req.Limit),
// 		Offset: int(req.Offset),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	for e := range entry {
// 		if err := stream.Send(&mediapb.ListMediaResponse{
// 			Endpoint: e.Endpoint,
// 			Name:     e.Name,
// 		}); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
