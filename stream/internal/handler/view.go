package handler

import (
	"fmt"
	"html/template"
	"io"

	"github.com/odit-bit/sone/stream/internal/domain"
)

// ////////////////////////////////////////////

const (
	_videoPlayer   = "VideoPlayer"
	_videoListPage = "listPage"
)

type ViewHandler struct {
	t map[string]*template.Template
}

func UI() *ViewHandler {
	t := make(map[string]*template.Template)
	t[_videoPlayer] = template.Must(template.New(_videoPlayer).Parse(VideoPlayerHTML))
	t[_videoListPage] = template.Must(template.New(_videoListPage).Parse(VideoListHTML))

	return &ViewHandler{
		t: t,
	}
}

func (h *ViewHandler) PlaybackPage(w io.Writer, data VideoPlayerData) error {
	tmpl, ok := h.t[_videoPlayer]
	if !ok {
		return fmt.Errorf("template not exist")
	}

	return tmpl.Execute(w, data)
}

func (h *ViewHandler) VideoListPage(w io.Writer, data []domain.Entry) error {
	tmpl, ok := h.t[_videoListPage]
	if !ok {
		return fmt.Errorf("template not exist")
	}

	return tmpl.Execute(w, data)
}
