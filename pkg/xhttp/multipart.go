package xhttp

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
)

type formData struct {
	ContentLength int64
	ContentType   string
	Filename      string
	Body          io.ReadCloser
}

func (f *formData) Close() error {
	return f.Body.Close()
}

func HandleMultipart(req *http.Request, formName string) (*formData, error) {
	length, err := strconv.Atoi(req.Header.Get("content-length"))
	if err != nil {
		return nil, err
	}

	mt, param, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if mt != "multipart/form-data" {
		return nil, fmt.Errorf("not mulitpart request")
	}

	reader := multipart.NewReader(req.Body, param["boundary"])
	part, err := reader.NextPart()
	if err != nil {
		return nil, err
	}

	ct := part.Header.Get("content-type")
	cd, param, err := mime.ParseMediaType(part.Header.Get("content-disposition"))
	if err != nil {
		return nil, err
	}
	if cd != "form-data" {
		return nil, fmt.Errorf("not form-data request")
	}

	name := param["name"]
	if name != formName {
		return nil, fmt.Errorf("wrong form name, got '%s' expect 'file' ", name)
	}
	filename := param["filename"]

	return &formData{
		ContentLength: int64(length),
		ContentType:   ct,
		Filename:      filename,
		Body:          part,
	}, nil
}
