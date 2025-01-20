package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/grafov/m3u8"
	"github.com/odit-bit/sone/media/internal/application"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger *logrus.Logger
	app    *application.Media
}

func NewHLSHandler(app *application.Media, logger *logrus.Logger) *handler {
	logrus.New()
	return &handler{
		logger: logger,
		app:    app,
	}
}

func (h *handler) Register(mux *chi.Mux) {
	mux.Put("/media/{id}/{segment}", h.HandlePutSegment)
	mux.Get("/media/{id}/{segment}", h.HandleGetSegment)
}

func (h *handler) HandleGetSegment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	segment := chi.URLParam(r, "segment")

	key := filepath.Join(id, segment)
	seg, err := h.app.Get(r.Context(), application.GetSegmentArgs{Name: key})
	if err != nil {
		h.logger.Error("media service GET:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	defer seg.Body.Close()
	if _, err := io.Copy(w, seg.Body); err != nil {
		h.logger.Error("media service GET:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	
}

func (h *handler) HandlePutSegment(w http.ResponseWriter, r *http.Request) {

	// header := r.Header.Get("X-Ingress-Forwarded-For")
	// h.logger.Debug("media service:", r.Header)

	id := chi.URLParam(r, "id")
	segment := chi.URLParam(r, "segment")

	ext := filepath.Ext(segment)
	key := filepath.Join(id, segment)

	switch ext {
	case ".m3u8":
		// read body
		defer r.Body.Close()

		pl, err := parseM3u8(r.Body)
		if err != nil {
			h.logger.Debug("failed parse .m3u8 ", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		defer pl.Close()
		if err := h.app.Put(r.Context(), application.PutSegmentArgs{
			Key:  key,
			Body: pl.Encode(),
		}); err != nil {
			h.logger.Error("media api streamHLS handler failed save playlist ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

	case ".ts":

		// read body
		defer r.Body.Close()

		// path := fmt.Sprintf("%s/%s", id, segment)
		if err := h.app.Put(r.Context(), application.PutSegmentArgs{
			Key:  key,
			Body: r.Body,
		}); err != nil {
			h.logger.Error("media api streamHLS handler failed save segment ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

	default:
		h.logger.Debug("invalid extention")
		http.Error(w, fmt.Sprintf("invalid extention %v", ext), http.StatusBadRequest)
		return
	}

}

func parseM3u8(r io.Reader) (*m3u8.MediaPlaylist, error) {
	pl, t, err := m3u8.DecodeFrom(r, true)
	if err != nil {
		return nil, fmt.Errorf("error parsing .m3u8: %v", err)
	}
	switch t {
	case m3u8.MASTER:
		// master := pl.(*m3u8.MasterPlaylist)
		// fmt.Println("got master playlist", master)
		return nil, fmt.Errorf("not support master playlist, only send media playlist")

	case m3u8.MEDIA:
		media := pl.(*m3u8.MediaPlaylist)
		// fmt.Println("target_dur", media.TargetDuration, "ext-x-media-sequence", media.SeqNo, "lastSequence", media.Segments[media.Count()])
		return media, nil
	default:
		return nil, fmt.Errorf("unknown playlist type %T. this is a bug", t)
	}

}
