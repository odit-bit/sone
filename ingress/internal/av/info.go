package av

import (
	"encoding/json"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type info struct {
	Index   int    `json:"index"`
	Codec   string `json:"codec_name"`
	W       int    `json:"width"`
	H       int    `json:"height"`
	Fps     int    `json:"r_frame_rate"`
	BitRate int    `json:"bit_rate"`
}

func Info(filename string) (*info, error) {
	v, err := ffmpeg_go.Probe(filename)
	if err != nil {
		return nil, err
	}
	var i info
	if err := json.Unmarshal([]byte(v), &i); err != nil {
		return nil, err
	}

	return &i, nil

}
