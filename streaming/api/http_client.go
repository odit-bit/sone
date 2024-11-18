package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/odit-bit/sone/streaming/internal/domain"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

// type SegmentQuery struct {
// 	StreamKey  string
// 	StreamName string
// 	Segment    string
// }

// func (c *Client) GetSegment(ctx context.Context, query SegmentQuery) (*domain.Segment, error) {
// 	path := "file"
// 	u, err := url.Parse(fmt.Sprintf("http://%v/%v", c.addr, path))
// 	if err != nil {
// 		log.Fatalf("streaming-client-api, invalid address %v error: %v", c.addr, err)
// 	}
// 	u.Query().Add("key", query.StreamKey)
// 	u.Query().Add("Name", query.StreamName)
// 	u.Query().Add("segment", query.Segment)

// 	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
// 	if err != nil {
// 		return nil, errors.Join(err, fmt.Errorf("streaming-client-api, url %v error: %v", u.String(), err))
// 	}

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		// log.Fatalf("streaming-client-api, invalid address %v error: %v", c.addr, err)
// 		return nil, err
// 	}

// 	if res.StatusCode != 200 && res.StatusCode != 202 {
// 		return nil, errors.Join(err, fmt.Errorf("failed get entries, got status code %v", res.StatusCode))
// 	}
// }

type ListQuery struct {
	Limit int
}

// List implements domain.EntryRepository.
func (c *Client) List(ctx context.Context, query ListQuery) (<-chan domain.Entry, error) {
	path := "entries"
	u, err := url.Parse(fmt.Sprintf("http://%v/%v", c.addr, path))
	if err != nil {
		log.Fatalf("streaming-client-api, invalid address %v error: %v", c.addr, err)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("streaming-client-api, url %v error: %v", u.String(), err))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// log.Fatalf("streaming-client-api, invalid address %v error: %v", c.addr, err)
		return nil, err
	}

	if res.StatusCode != 200 && res.StatusCode != 202 {
		return nil, errors.Join(err, fmt.Errorf("failed get entries, got status code %v", res.StatusCode))
	}

	eC := make(chan domain.Entry)

	go func() {
		defer res.Body.Close()
		var entries []domain.Entry
		json.NewDecoder(res.Body).Decode(&entries)
		for _, e := range entries {
			eC <- e
		}
		close(eC)
	}()
	return eC, nil
}
