package grpc

// // CreateStream implements userpb.StreamServiceServer.
// func (s *Server) CreateStream(ctx context.Context, in *userpb.CreateStreamRequest) (*userpb.StreamInfo, error) {
// 	resp, err := s.app.InsertStream(ctx, application.CreateStreamArgs{Token: in.GetToken(), Title: in.GetToken()})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &userpb.StreamInfo{
// 		Key:             resp.Key,
// 		Id:              resp.ID,
// 		Title:           resp.Title,
// 		Created:         timestamppb.New(resp.CreatedAt),
// 		LastStartStream: timestamppb.New(resp.LastStartStream),
// 		LastEndStream:   timestamppb.New(resp.LastEndStream),
// 		IsStreaming:     resp.IsStreaming,
// 	}, nil
// }

// func (s *Server) StartStream(ctx context.Context, in *userpb.StartStreamRequest) (*userpb.StartStreamResponse, error) {
// 	err := s.app.StartStream(ctx, application.StartStreamArgs{
// 		StreamKey: in.StreamKey,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &userpb.StartStreamResponse{ErrorStatus: 0, ErrorMessage: ""}, nil
// }

// func (s *Server) EndStream(ctx context.Context, in *userpb.EndStreamRequest) (*userpb.EndStreamResponse, error) {
// 	err := s.app.EndStream(ctx, application.EndStreamArgs{
// 		StreamKey: in.StreamKey,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &userpb.EndStreamResponse{}, nil
// }

// func (s *Server) GetStream(ctx context.Context, in *userpb.GetStreamRequest) (*userpb.StreamInfo, error) {
// 	val, ok, err := s.app.GetStream(ctx, in.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !ok {
// 		return nil, fmt.Errorf("stream not found")
// 	}

// 	return &userpb.StreamInfo{
// 		Id:              val.ID,
// 		Title:           val.Title,
// 		Key:             val.Key,
// 		Created:         timestamppb.New(val.CreatedAt),
// 		IsStreaming:     val.IsStreaming,
// 		LastStartStream: timestamppb.New(val.LastStartStream),
// 	}, nil
// }

// func (s *Server) ListStream(in *userpb.ListRequest, resp userpb.StreamService_ListStreamServer) error {
// 	c, err := s.app.ListStream(resp.Context(), 1000, 0)
// 	if err != nil {
// 		return err
// 	}

// 	for val := range c {
// 		if err := resp.Send(&userpb.StreamInfo{
// 			Id:              val.ID,
// 			Title:           val.Title,
// 			Key:             val.Key,
// 			Created:         timestamppb.New(val.CreatedAt),
// 			IsStreaming:     val.IsStreaming,
// 			LastStartStream: timestamppb.New(val.LastStartStream),
// 		}); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
