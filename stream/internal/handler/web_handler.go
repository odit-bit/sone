package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/odit-bit/sone/stream/internal/domain"
	"github.com/sirupsen/logrus"
)

type WebHandler struct {
	fh   domain.EntryRepository
	view *ViewHandler
	// addr   string
	logger  *logrus.Logger
	keyRepo domain.StreamKeyRepo
}

func NewWebHandler(fh domain.EntryRepository, vh *ViewHandler, kr domain.StreamKeyRepo) *WebHandler {
	h := &WebHandler{
		fh:      fh,
		view:    vh,
		keyRepo: kr,

		logger: logrus.StandardLogger(),
	}
	return h
}

// start is non blocking, dont forget invoke close when done
func (s *WebHandler) Handler() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Logger) //
	r.Use(middleware.Recoverer)

	// html page
	r.Get("/", s.indexPage)
	r.Get("/stream", s.createStreamPage)
	r.Get("/view/{stream}", s.playbackPage)

	// api callback by the html page
	r.Get("/file/{name}/{segment}", s.segment)
	r.Get("/stream/create", s.createStream)

	return r

}

func (h *WebHandler) indexPage(w http.ResponseWriter, r *http.Request) {
	list := h.fh.List(0, 0)

	resp := []domain.Entry{}
	for e := range list {
		e.Path = path.Join("view", filepath.Base(e.Name))
		resp = append(resp, e)
	}

	err := h.view.VideoListPage(w, resp)
	if err != nil {
		log.Println(err)
	}
}

func (h *WebHandler) playbackPage(w http.ResponseWriter, r *http.Request) {
	stream := chi.URLParam(r, "stream")
	source := filepath.Join("/file", stream, "index.m3u8")
	h.view.PlaybackPage(w, VideoPlayerData{SourceURL: source})
}

func (h *WebHandler) segment(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	segment := chi.URLParam(r, "segment")

	path := filepath.Join(name, segment)
	seg := domain.Segment{}
	if err := h.fh.GetSegment(path, &seg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer seg.Body.Close()
	_, err := io.Copy(w, seg.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) createStreamPage(w http.ResponseWriter, r *http.Request) {
	if err := h.view.StreamKeyPage(w, "stream/create"); err != nil {
		h.logger.Error(fmt.Sprintf("failed to send 'create-stream-page-html' %v, this is a BUG ", err))
	}
}

func (h *WebHandler) createStream(w http.ResponseWriter, r *http.Request) {
	key := domain.CreateKey()
	h.keyRepo.Save(r.Context(), key)
	w.Write([]byte(string(key)))
}
