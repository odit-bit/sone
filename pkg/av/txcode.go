package av

import (
	"context"
	"fmt"
	"io"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	HLS_VOD   hlsPlaylist = "vod"
	HLS_EVENT hlsPlaylist = "event"
	HLS_lIVE  hlsPlaylist = "linear"
)

type hlsPlaylist string

type HLSArgs struct {
	Types hlsPlaylist
}

func PostHLS(ctx context.Context, uri string, flvReader io.Reader, param HLSArgs) error {
	//ffmpeg -re -i test-stream.mp4 -codec copy -f hls  -hls_list_size 0 http://localhost:6969/{stream-key}/{playlist}.m3u8
	// ffmpeg instance

	var (
		logErrOpt ffmpeg.KwArgs = ffmpeg.KwArgs{
			"hide_banner": "",
			"loglevel":    "error",
		}

		transCode ffmpeg.KwArgs = ffmpeg.KwArgs{
			"c:v": "copy",
			"c:a": "copy",
		}

		hlsOpt ffmpeg.KwArgs = ffmpeg.KwArgs{
			"f":                "hls",
			"hls_init_time":    5,
			"hls_time":         6,
			"hls_list_size":    5,
			"hls_segment_type": "mpegts",
			"hls_flags":        "program_date_time",

			"method":          "PUT",
			"http_user_agent": "UserAgent: 'sone/1.0 (compatible; linux)'",
			"headers":         fmt.Sprintf("X-Ingress-Forwarded-For: %s", uri),
		}
	)

	//validate
	switch param.Types {
	case HLS_EVENT:
	case HLS_VOD:
	case HLS_lIVE:
		// //video
		// transCode["c:v"] = "libx264"
		// transCode["b:v"] = "2800k"
		// transCode["maxrate:v"] = "2996k"
		// transCode["bufsize:v"] = "4200k"
		// //audio
		// transCode["c:a"] = "aac"
		// transCode["b:a"] = "160k"
		// transCode["ac"] = 2

	}

	var inArgs = []ffmpeg.KwArgs{{"f": "flv"}, logErrOpt}
	var outArgs = []ffmpeg.KwArgs{transCode, hlsOpt}

	s := ffmpeg.Input("pipe:", ffmpeg.MergeKwArgs(inArgs))
	s.Context = ctx
	s = s.Output(uri, ffmpeg.MergeKwArgs(outArgs)).OverWriteOutput().ErrorToStdOut()
	s = s.WithInput(flvReader)

	err := s.Run()
	if err != nil {
		return err
	}
	return nil
}
