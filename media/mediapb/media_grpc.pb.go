// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: media/mediapb/media.proto

// command to generate
// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./media/mediapb/media.proto

package mediapb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MediaService_GetSegment_FullMethodName = "/mediapb.MediaService/GetSegment"
)

// MediaServiceClient is the client API for MediaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MediaServiceClient interface {
	GetSegment(ctx context.Context, in *SegmentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SegmentResponse], error)
}

type mediaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMediaServiceClient(cc grpc.ClientConnInterface) MediaServiceClient {
	return &mediaServiceClient{cc}
}

func (c *mediaServiceClient) GetSegment(ctx context.Context, in *SegmentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SegmentResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MediaService_ServiceDesc.Streams[0], MediaService_GetSegment_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SegmentRequest, SegmentResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MediaService_GetSegmentClient = grpc.ServerStreamingClient[SegmentResponse]

// MediaServiceServer is the server API for MediaService service.
// All implementations must embed UnimplementedMediaServiceServer
// for forward compatibility.
type MediaServiceServer interface {
	GetSegment(*SegmentRequest, grpc.ServerStreamingServer[SegmentResponse]) error
	mustEmbedUnimplementedMediaServiceServer()
}

// UnimplementedMediaServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMediaServiceServer struct{}

func (UnimplementedMediaServiceServer) GetSegment(*SegmentRequest, grpc.ServerStreamingServer[SegmentResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GetSegment not implemented")
}
func (UnimplementedMediaServiceServer) mustEmbedUnimplementedMediaServiceServer() {}
func (UnimplementedMediaServiceServer) testEmbeddedByValue()                      {}

// UnsafeMediaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MediaServiceServer will
// result in compilation errors.
type UnsafeMediaServiceServer interface {
	mustEmbedUnimplementedMediaServiceServer()
}

func RegisterMediaServiceServer(s grpc.ServiceRegistrar, srv MediaServiceServer) {
	// If the following call pancis, it indicates UnimplementedMediaServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MediaService_ServiceDesc, srv)
}

func _MediaService_GetSegment_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SegmentRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MediaServiceServer).GetSegment(m, &grpc.GenericServerStream[SegmentRequest, SegmentResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MediaService_GetSegmentServer = grpc.ServerStreamingServer[SegmentResponse]

// MediaService_ServiceDesc is the grpc.ServiceDesc for MediaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MediaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mediapb.MediaService",
	HandlerType: (*MediaServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSegment",
			Handler:       _MediaService_GetSegment_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "media/mediapb/media.proto",
}
