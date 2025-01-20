package web

import "net/http"

type Router interface {
	Get(path string, fn http.HandlerFunc)
	Post(path string, fn http.HandlerFunc)
}
