package mediahttp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

type ClientGetResponse struct {
	Body io.ReadCloser
}

var ErrInvalidEndpoint = errors.New("invalid endpoint")

func (c *Client) Get(ctx context.Context, key string) (*ClientGetResponse, error) {
	uri := fmt.Sprintf("%s/media/%s", c.endpoint, key)
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 400, 500:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("media server response: %s", string(b))
	case 404:
		return nil, ErrInvalidEndpoint
	}

	return &ClientGetResponse{Body: resp.Body}, nil
}

type ClientPutResponse struct {
	IsSuccess bool
}

func (c *Client) Put(ctx context.Context, key string, r io.Reader) (*ClientPutResponse, error) {
	uri := fmt.Sprintf("%s/media/%s", c.endpoint, key)
	req, err := http.NewRequestWithContext(ctx, "PUT", uri, r)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 400, 404, 500:
		return nil, c.handleErrStatus(resp)
	}

	return &ClientPutResponse{IsSuccess: true}, nil
}

func (c *Client) handleErrStatus(res *http.Response) error {
	switch res.StatusCode {
	case 400, 500:
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("media server response: %s", string(b))
	case 404:
		return ErrInvalidEndpoint
	}

	return nil
}
