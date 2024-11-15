package main

import (
	"io"
	"log"
	"os"

	"github.com/bluenviron/gohlslib/v2"
	"github.com/bluenviron/gohlslib/v2/pkg/codecs"
	"github.com/bluenviron/mediacommon/pkg/formats/mpegts"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
)

func main() {
	videoTrack := &gohlslib.Track{
		Codec:     &codecs.H264{},
		ClockRate: 90000,
	}
	mux := gohlslib.Muxer{
		Tracks: []*gohlslib.Track{videoTrack},
	}

	mux.Start()

	f, err := os.Open("/mnt/d/wsl/test.flv")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	dec, err := flv.NewDecoder(f)
	if err != nil {
		log.Println("err")
		return
	}

	var flvTag flvtag.FlvTag
	for {
		//tag
		if err := dec.Decode(&flvTag); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Failed to decode: %+v", err)
		}

		switch flvTag.TagType {
		case flvtag.TagTypeAudio:
			log.Printf("TagAUDIO: %+v", flvTag)
			// mux.WriteH264(videoTrack, flvTag.Timestamp)

		case flvtag.TagTypeVideo:
			log.Printf("TagVIDEO: %+v", flvTag)

		case flvtag.TagTypeScriptData:
			log.Printf("TagSCRIPTDATA: %+v", flvTag)

		}
		flvTag.Close() // Discard unread buffers
	}

}

func findH264Track(r *mpegts.Reader) *mpegts.Track {
	for _, track := range r.Tracks() {
		if _, ok := track.Codec.(*mpegts.CodecH264); ok {
			return track
		}
	}
	return nil
}
