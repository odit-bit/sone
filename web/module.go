package web

import (
	"context"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	streamingapi "github.com/odit-bit/sone/streaming/api"
)

func Register(addr string, view *ViewHandler, r *chi.Mux) {
	s := handler{
		view:    view,
		entries: streamingapi.NewClient(addr),
	}
	r.Get("/", s.indexPage)
	r.Get("/stream", s.createStreamPage)
	r.Get("/view/{stream-key}/{stream-name}", s.playbackPage)
}

type handler struct {
	view    *ViewHandler
	entries *streamingapi.Client
}

func (h *handler) indexPage(w http.ResponseWriter, r *http.Request) {
	list, err := h.entries.List(context.Background(), streamingapi.ListQuery{})
	if err != nil {
		log.Println(err)
		return
	}

	resp := []Entry{}
	for e := range list {
		resp = append(resp, Entry{
			Endpoint: e.Endpoint,
			Name:     e.Name,
		})
	}

	err = h.view.VideoListPage(w, resp)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) playbackPage(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "stream-key")
	name := chi.URLParam(r, "stream-name")
	source := filepath.Join("/file", key, name, "index.m3u8")
	h.view.PlaybackPage(w, VideoPlayerData{SourceURL: source})
}

func (h *handler) createStreamPage(w http.ResponseWriter, r *http.Request) {
	if err := h.view.StreamKeyPage(w, "stream/create"); err != nil {
		// h.logger.Error(fmt.Sprintf("failed to send 'create-stream-page-html' %v, this is a BUG ", err))
	}
}

func StartUp(listener net.Listener, mux *chi.Mux) error {
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return err
	}
	if host == "" || host == "::" {
		host = "localhost"
	}
	Register(net.JoinHostPort(host, port), UI(), mux)
	return err
}
