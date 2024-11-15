package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/livekit/go-rtmp"
	"github.com/livekit/go-rtmp/message"
	"github.com/odit-bit/sone/pkg/buffer"
	"github.com/odit-bit/sone/rtmp/internal/av"
	"github.com/odit-bit/sone/rtmp/internal/domain"
	"github.com/yutopp/go-flv"
	"github.com/yutopp/go-flv/tag"
)

const (
	DefaultMaxBitrate uint32 = 5000
)

type CreateHLSHandler struct {
	rtmp.DefaultHandler
	repo domain.HLSRepository

	ctx    context.Context
	cancel func()

	encoder *flv.Encoder

	pr   *io.PipeReader
	pw   *io.PipeWriter
	pErr chan error
	// tcErr error

	maxBitRate uint32
	limiter    *av.LimitedBitrate
}

func NewCreateHLSHandler(ctx context.Context, hlsRepo domain.HLSRepository, opts ...CreateHLSOpt) *CreateHLSHandler {
	pr, pw := io.Pipe()
	ictx, cancel := context.WithCancel(ctx)
	h := &CreateHLSHandler{
		DefaultHandler: rtmp.DefaultHandler{},
		repo:           hlsRepo,
		ctx:            ictx,
		cancel:         cancel,
		encoder:        nil,
		pr:             pr,
		pw:             pw,
		pErr:           make(chan error, 1),
		maxBitRate:     DefaultMaxBitrate,
		limiter:        nil,
	}

	for _, opt := range opts {
		opt(h)
	}
	return h

}

func (h *CreateHLSHandler) OnServe(conn *rtmp.Conn) {}

func (h *CreateHLSHandler) OnConnect(timestamp uint32, cmd *message.NetConnectionConnect) error {
	// sURL, err := url.Parse(cmd.Command.TCURL)
	// if err != nil {
	// 	return fmt.Errorf("failed parse url from stream")
	// }
	// p, ok := sURL.User.Password()
	// if !ok || p != "password" {
	// 	return fmt.Errorf("unauthorized")
	// }

	// log.Printf("OnConnect-stream: %#v", cmd)
	log.Println("Connect STREAM")
	return nil
}

func (h *CreateHLSHandler) OnCreateStream(timestamp uint32, cmd *message.NetConnectionCreateStream) error {
	// log.Printf("OnCreateStream: %#v", cmd)
	return nil
}

// client publish to stream command
func (h *CreateHLSHandler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *message.NetStreamPublish) error {
	// log.Printf("OnPublish: %#v", cmd)
	// (example) Reject a connection when PublishingName is empty
	if cmd.PublishingName == "" {
		return errors.New("PublishingName is empty")
	}

	// //setup limiter
	// limiter := av.NewBitrateRejectorReader(h.pr, h.maxBitRate)
	// h.limiter = limiter

	// setup ingestion
	// it should ingest the serialized rtmp payload (flv)
	path := strings.Split(cmd.PublishingName, "/")
	if len(path) != 2 {
		return fmt.Errorf("invalid stream key/name ")
	}
	go func() {
		// log.Println("RUNNING Trancoder")
		c := h.repo.TranscodeFLV(h.ctx, path[0], path[1], h.pr)
		select {
		case <-h.ctx.Done():
		case resp := <-c:
			if resp != nil {
				h.pErr <- resp
			}
		}
		h.pr.Close()
		// log.Printf("ingress transcoder finish: %v \n", resp)
		// log.Println("EXIT Trancoder")
	}()

	//setup encoder to flv
	enc, err := flv.NewEncoder(h.pw, flv.FlagsAudio|flv.FlagsVideo)
	if err != nil {
		return h.pipeError(errors.Join(err, fmt.Errorf("failed created flv encoder")))
	}

	h.encoder = enc
	return nil
}

// client play from stream command
func (h *CreateHLSHandler) OnPlay(ctx *rtmp.StreamContext, timestamp uint32, cmd *message.NetStreamPlay) error {
	// log.Printf("OnPlay-stream: %#v", cmd)
	return nil
}

func (h *CreateHLSHandler) OnSetDataFrame(timestamp uint32, data *message.NetStreamSetDataFrame) error {

	r := buffer.Pool.Get()
	defer buffer.Pool.Put(r)

	var script tag.ScriptData
	if err := tag.DecodeScriptData(r, &script); err != nil {
		// log.Printf("failed to decode script data , err = %+v", err)
		return h.pipeError(err)
	}

	// log.Printf("set data frame: Script= %#v", script)

	if err := h.encoder.Encode(&tag.FlvTag{
		TagType:   tag.TagTypeScriptData,
		Timestamp: timestamp,
		Data:      &script,
	}); err != nil {
		// log.Printf("failed to write script data, err : %+v", err)
		return h.pipeError(err)
	}

	return nil
}

func (h *CreateHLSHandler) OnAudio(timestamp uint32, payload io.Reader) error {

	var audio tag.AudioData
	if err := tag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}

	body := buffer.Pool.Get()
	defer buffer.Pool.Put(body)

	if _, err := io.Copy(body, audio.Data); err != nil {
		return err
	}

	audio.Data = body

	if err := h.encoder.Encode(&tag.FlvTag{
		TagType:   tag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      &audio,
	}); err != nil {
		// log.Printf("Failed to write audio: Err = %+v", err)
		return h.pipeError(err)
	}

	return nil
}

func (h *CreateHLSHandler) OnVideo(timestamp uint32, payload io.Reader) error {

	var video tag.VideoData

	if err := tag.DecodeVideoData(payload, &video); err != nil {
		return fmt.Errorf("failed to decode video, err : %v", err)
	}

	body := bytes.Buffer{}
	if _, err := io.Copy(&body, video.Data); err != nil {
		return h.pipeError(err)
	}
	video.Data = &body

	if err := h.encoder.Encode(&tag.FlvTag{
		TagType:   tag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &video,
	}); err != nil {
		return h.pipeError(fmt.Errorf("failed to encode video, err; %v", err))
	}

	return nil
}

func (h *CreateHLSHandler) OnClose() {

	if h.pw != nil {
		h.pw.Close()
	}

	if h.pr != nil {
		h.pr.Close()
	}

	if h.limiter != nil {
		h.limiter.Close()
	}

	if h.cancel != nil {
		h.cancel()
	}
}

// to check if error exist
func (h *CreateHLSHandler) pipeError(err error) error {
	select {
	case pErr, ok := <-h.pErr:
		if ok {
			return pErr
		}
		if pErr != nil {
			err = errors.Join(err, pErr)
		}
	default:
	}
	return err

}

type CreateHLSOpt func(h *CreateHLSHandler)

func WithBitRateLimit(kbps uint32) CreateHLSOpt {
	if kbps < 1000 {
		log.Fatal("minimum value of bitrate is 1000 kbps")
	}

	return func(h *CreateHLSHandler) {
		h.maxBitRate = kbps
	}
}
