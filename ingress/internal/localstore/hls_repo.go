package localstore

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/odit-bit/sone/ingress/internal/domain"
	"github.com/odit-bit/sone/ingress/internal/rpc"
	"github.com/odit-bit/sone/pkg/av"
)

var _ domain.HLSRepository = (*hlsRepository)(nil)

type hlsRepository struct {
	streams *rpc.Streams
}

func NewHLSRepository(streamsAPI *rpc.Streams) *hlsRepository {
	// afs := afero.NewBasePathFs(afero.NewOsFs(), root)
	return &hlsRepository{
		streams: streamsAPI,
	}
}

// TranscodeFLV implements domain.HLSRepository.
func (h *hlsRepository) TranscodeFLV(ctx context.Context, stream_key string, r io.Reader) <-chan *domain.OpError {
	respC := make(chan *domain.OpError, 1)

	//START STREAMING
	resp, err := h.streams.StartStream(ctx, stream_key)
	if err != nil {
		log.Println("this is BUG it should show on main error Log")
		respC <- domain.NewOpError(err.Error())
		close(respC)
		return respC
	}

	writeErr := func(err error) {
		opErr := domain.NewOpError(err.Error())
		respC <- opErr
		close(respC)
	}

	go func() {

		//STOP STREAMING
		defer func() {
			eCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := h.streams.EndStream(eCtx, stream_key)
			if err != nil {
				log.Printf("ingress failed send endstream event, status: %v \n", err)
				// writeErr(err)
			}
		}()

		// err := h.postHLS(ctx, fmt.Sprintf("%v/%v/%v", "http://localhost:9797/media", resp.ID, "index.m3u8"), r)
		err := av.PostHLS(ctx, fmt.Sprintf("%v/%v/%v", "http://localhost:9797/media", resp.ID, "index.m3u8"), r, av.HLSArgs{Types: av.HLS_lIVE})

		if err != nil {
			// err = errors.Join(err, h.root.RemoveAll(dir))
			writeErr(err)
			return
		}

	}()

	return respC

}

// func (h *hlsRepository) postHLS(ctx context.Context, uri string, flvReader io.Reader) error {
// 	//ffmpeg -re -i test-stream.mp4 -codec copy -f hls  -hls_list_size 0 http://localhost:6969/{stream-key}/{playlist}.m3u8
// 	// ffmpeg instance

// 	// ffmpeg args

// 	var inArgs = []ffmpeg.KwArgs{{"f": "flv"}, av.LogErrOpt}
// 	var outArgs = []ffmpeg.KwArgs{av.TransCode, {

// 		"f":                "hls",
// 		"hls_init_time":    6,
// 		"hls_time":         10,
// 		"hls_list_size":    5,
// 		"hls_segment_type": "mpegts",
// 	}}

// 	s := ffmpeg.Input("pipe:", inArgs...)
// 	s.Context = ctx
// 	s = s.Output(uri, outArgs...).OverWriteOutput().ErrorToStdOut()
// 	s = s.WithInput(flvReader)

// 	err := s.Run()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
