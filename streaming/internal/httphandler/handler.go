package httphandler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/odit-bit/sone/streaming/internal/application"
	"github.com/odit-bit/sone/streaming/internal/application/commands"
	"github.com/odit-bit/sone/streaming/internal/application/queries"
	"github.com/odit-bit/sone/streaming/internal/domain"
)

func Register(app application.App, mux *chi.Mux) {
	h := handler{app: app}
	mux.Get("/file/{key}/{name}/{segment}", h.handleGetSegment)
	mux.Get("/entries", h.handleGetEntries)
	mux.Get("/stream/create", h.handleCreateStream)

}

type handler struct {
	app application.App
}

func (h *handler) handleCreateStream(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateStreamKeyCMD
	key, err := h.app.CreateStreamKey(r.Context(), cmd)
	_ = err
	w.Write([]byte(key))
}

func (h *handler) handleGetSegment(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	name := chi.URLParam(r, "name")
	reqSegment := chi.URLParam(r, "segment")

	path := filepath.Join(key, name, reqSegment)
	seg := queries.GetSegment{
		Name: path,
	}
	segment, err := h.app.GetSegment(r.Context(), seg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer segment.Body.Close()
	_, err = io.Copy(w, segment.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) handleGetEntries(w http.ResponseWriter, r *http.Request) {
	list, _ := h.app.GetEntries(r.Context(), queries.GetEntry{Limit: 0, Offset: 0})

	resp := []domain.Entry{}
	for e := range list {
		e.Endpoint = filepath.Join("/view", e.Endpoint)
		resp = append(resp, e)
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		_, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", filename, line, err)
		// h.logger.Errorf("%v|%v handleGetEntries failed encode to javascript", filename, line)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

}
