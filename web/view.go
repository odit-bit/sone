package web

import (
	"fmt"
	"html/template"
	"io"
)

// ////////////////////////////////////////////

const (
	_videoPlayer         = "VideoPlayer"
	_videoListPage       = "listPage"
	_createStreamKeyPage = "createStreamKeyPage"
)

type ViewHandler struct {
	t map[string]*template.Template
}

func UI() *ViewHandler {
	t := make(map[string]*template.Template)
	t[_videoPlayer] = template.Must(template.New(_videoPlayer).Parse(VideoPlayerHTML))
	t[_videoListPage] = template.Must(template.New(_videoListPage).Parse(VideoListHTML))
	t[_createStreamKeyPage] = template.Must(template.New(_createStreamKeyPage).Parse(StreamKeyCreateHTML))

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

func (h *ViewHandler) VideoListPage(w io.Writer, data []Entry) error {
	tmpl, ok := h.t[_videoListPage]
	if !ok {
		return fmt.Errorf("template not exist")
	}

	return tmpl.Execute(w, data)
}

func (h *ViewHandler) StreamKeyPage(w io.Writer, endpointCB string) error {
	tmpl, ok := h.t[_createStreamKeyPage]
	if !ok {
		return fmt.Errorf("template not exist, %v ", _createStreamKeyPage)
	}

	agg := map[string]string{"CreateStreamKeyURL": endpointCB}
	return tmpl.Execute(w, agg)
}
