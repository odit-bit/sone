package web

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/klauspost/compress/gzhttp"
	"github.com/odit-bit/sone/media/mediahttp"
	"github.com/odit-bit/sone/streaming/streamingpb"
	"github.com/odit-bit/sone/web/internal/session"
)

//////////////////////////////////////////

type StreamHandler struct {
	rootPage     Renderer
	createPage   Renderer
	playbackPage Renderer

	streamAPI streamingpb.LiveStreamClient
	mediaAPI  *mediahttp.Client

	sm session.Manager
}

func NewStreamHandler(sm session.Manager, streamAPI streamingpb.LiveStreamClient, mediaBackend *mediahttp.Client) *StreamHandler {
	return &StreamHandler{
		rootPage:     NewStreamRootTemplate(),
		playbackPage: NewPlaybackTemplate(),
		createPage:   NewCreateStreamTemplate(""),
		streamAPI:    streamAPI,
		mediaAPI:     mediaBackend,
		sm:           sm,
	}
}

func (h *StreamHandler) Handle(path string, mux *chi.Mux) {
	//gzip middleware
	// root
	mux.HandleFunc(path, h.renderRootPage("/stream/playback", "/stream/create"))
	mux.Post("/stream/create", h.handleGenerateStreamKey())

	// playback
	mux.Get("/stream/playback/{stream-key}", h.renderPlayback("/stream"))
	mux.Get("/stream/{id}/{segment}", gzhttp.GzipHandler(h.handlePlayback()))

}

func (h *StreamHandler) renderRootPage(playbackPrefix, createStreamCB string) http.HandlerFunc {
	playbackUri := func(endpoint string) string {
		return path.Join(playbackPrefix, endpoint)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := h.streamAPI.List(context.Background(), &streamingpb.ListRequest{})
		if err != nil {
			log.Println("stream handler failed get stream list:", err)
			http.Error(w, "", http.StatusServiceUnavailable)
			return
		}

		list := []streamListArgs{}
		for {
			msg, err := c.Recv()
			if err != nil {
				log.Println("stream handler get list debug:", err)
				break
			}
			list = append(list, streamListArgs{
				PlaybackUrl: playbackUri(msg.Id),
				Title:       msg.Title,
			})
		}
		h.rootPage.Render(w, StreamRootArgs{List: list, CreateStreamCB: createStreamCB})
	}
}

func (h *StreamHandler) renderPlayback(playbackCb string) http.HandlerFunc {
	uri := func(key string) string {
		return filepath.Join("/", playbackCb, key, "index.m3u8")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "stream-key")

		if err := h.playbackPage.Render(w, PlaybackData{SourceURL: uri(key)}); err != nil {
			log.Println("stream handler failed render playback page:", err)
			http.Error(w, "", 500)
			return
		}
	}
}

func (h *StreamHandler) handlePlayback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		reqSegment := chi.URLParam(r, "segment")
		ext := path.Ext(reqSegment)

		switch ext {
		case ".ts":
			//video
			w.Header().Set("Content-Type", "video/mp2t")

		case ".m3u8":
			//hls-playlist
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

		}

		// segment, err := h.mediaAPI.GetSegment(r.Context(), &mediapb.SegmentRequest{
		// 	Path: fmt.Sprintf("%s/%s", id, reqSegment),
		// })
		segment, err := h.mediaAPI.Get(r.Context(), fmt.Sprintf("%s/%s", id, reqSegment))

		if err != nil {
			log.Println("web stream handler failed to fetch segment: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer segment.Body.Close()

		_, err = io.Copy(w, segment.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *StreamHandler) handleGenerateStreamKey() http.HandlerFunc {
	tmpl := template.Must(template.New("generateStreamKeyResponse").Parse("<p>your stream-key: {{.Key}}</p>"))

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.sm.Load(r)
		if err != nil {
			log.Println("create stream handler failed get session:", err)
			http.Error(w, "failed get session", http.StatusInternalServerError)
			return
		}

		if ok := data.IsLogin(); !ok {
			http.Error(w, "not login", http.StatusBadRequest)
			return
		}

		key := data.GetString("stream_key")

		if key == "" {
			resp, err := h.streamAPI.Insert(r.Context(), &streamingpb.InsertRequest{
				Token: data.GetString("token"),
				Title: "",
			})
			if err != nil {
				log.Println("create stream handler error:", err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}

			data.Put("stream_key", string(resp.Key))
			if err := h.sm.Save(data, w); err != nil {
				log.Println("create stream handler failed save session:", err)
				http.Error(w, "failed to save session", http.StatusInternalServerError)
				return
			}
			key = resp.Key
		}

		tmpl.Execute(w, map[string]string{"Key": key})
	}

}
