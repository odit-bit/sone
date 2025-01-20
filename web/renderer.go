package web

import (
	"io"
)

type Renderer interface {
	Render(w io.Writer, data any) error
}
