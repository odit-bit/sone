package av

import (
	"context"
	"io"
	"path/filepath"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

const (
	HLS_PLAYLIST = "index.m3u8"
)

var (
	LogErrOpt ffmpeg_go.KwArgs = ffmpeg_go.KwArgs{
		"hide_banner": "",
		"loglevel":    "error",
	}

	TransCode ffmpeg_go.KwArgs = ffmpeg_go.KwArgs{
		"c:v": "copy",
		"c:a": "copy",
	}

	HlsFileOpt ffmpeg_go.KwArgs = ffmpeg_go.KwArgs{
		"f":                "hls",
		"hls_time":         1,
		"hls_list_size":    5,
		"hls_segment_type": "mpegts",
		// "hls_playlist_type": "vod",
		"hls_flags": "delete_segments",

		// "f":                  "ssegment", //not type (stream_segment)
		// "segment_list_flags": "+live",
		// "segment_format":     "mpegts",
		// "segment_time":       1,
		// "segment_list_size":  6,
		// "segment_list_type":  "m3u8",
	}

	HLSStreamOpt ffmpeg_go.KwArgs = ffmpeg_go.KwArgs{
		"f":                    "hls",
		"hls_time":             1,
		"hls_list_size":        5,
		"hls_segment_type":     "mpegts",
		"hls_segment_filename": "pipe:1",
		"pipe:1.m3u8":          "",
	}

	// HlsFileOpt ffmpeg_go.KwArgs = ffmpeg_go.KwArgs{
	// 	"f":                "hls",
	// 	"hls_time":         5,
	// 	"hls_list_size":    0, //if not '0', playlist only store 6 of segment file
	// 	"hls_segment_type": "mpegts",
	// 	// "hls_playlist_type": "vod",
	// 	// "hls_flags":         "independent_segments",
	// }

	// TransMux ffmpeg_go.KwArgs = ffmpeg_go.KwArgs{
	// 	"c:v":       "libx264",
	// 	"c:a":       "aac",
	// 	"s":         "1280x720",
	// 	"sws_flags": "bilinear",
	// 	"r":         30,
	// 	"b:v":       "4400k",
	// 	// "bufsize ":  4400,
	// 	"b:a": "128k",

	// 	"profile:v": "main",
	// 	"preset":    "veryfast",
	// }
)

// func HlsStreamOptFunc(outpath string) ffmpeg_go.KwArgs {
// 	HlsStreamOpt["segment_list"] = filepath.Join(outpath, "index.m3u8")
// 	return HlsStreamOpt
// }

func ReadHLSFrom(ctx context.Context, w io.Writer, r io.Reader) error {
	args := []ffmpeg_go.KwArgs{}
	args = append(args, TransCode, HLSStreamOpt)

	// ffmpeg instance
	s := ffmpeg_go.Input("pipe:", ffmpeg_go.KwArgs{"f": "flv"}, LogErrOpt)
	s.Context = ctx
	s = s.Output("pipe:", ffmpeg_go.MergeKwArgs(args)).OverWriteOutput().ErrorToStdOut().
		WithOutput(w).
		WithInput(r)

	err := s.Run()
	if err != nil {
		return err
	}
	return nil
}

// will create segment file from reader and write to disk,
func WriteFLVToHLS(ctx context.Context, path string, r io.Reader) error {
	// ffmpeg -re -i <input-file.mp4> -c:a copy -c:v copy -f hls -hls_time 5 -hls_list_size 0 <folder/index.m3u8>

	// path := filepath.Join(dir, streamName)
	// os.MkdirAll(path, 0755)
	// indexPath, _ := filepath.Abs(filepath.Join(path, HLS_PLAYLIST))

	dir := filepath.Dir(path)
	segmentName := filepath.Join(dir, "segment%03d.ts")

	//ffmpeg args
	args := []ffmpeg_go.KwArgs{}
	args = append(args, TransCode, HlsFileOpt, ffmpeg_go.KwArgs{"hls_segment_filename": segmentName})

	// ffmpeg instance
	s := ffmpeg_go.Input("pipe:", ffmpeg_go.KwArgs{"f": "flv"}, LogErrOpt)
	s.Context = ctx
	s = s.Output(path, ffmpeg_go.MergeKwArgs(args)).OverWriteOutput().ErrorToStdOut()
	s = s.WithInput(r)

	err := s.Run()
	if err != nil {
		return err
	}
	return nil

}
