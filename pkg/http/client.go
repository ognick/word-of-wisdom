package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/ognick/word_of_wisdom/pkg/logger"
)

type handler func(ctx context.Context, status int, header http.Header, body []byte) error

type Client struct {
	addr    string
	method  string
	handler handler
	client  *http.Client
	log     logger.Logger
}

func NewClient(
	log logger.Logger,
	addr string,
	method string,
	handler handler,
) *Client {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	if host == "" {
		host = "http://0.0.0.0"
	}
	address := host + ":" + port
	return &Client{
		log:     log,
		addr:    address,
		method:  method,
		handler: handler,
		client:  &http.Client{},
	}
}

func (c *Client) Request(
	ctx context.Context,
	header http.Header,
	reqBody io.Reader,
) error {
	req, err := http.NewRequestWithContext(ctx, c.method, c.addr, reqBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	for key, values := range header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request:%w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response:%w", err)
	}

	return c.handler(ctx, resp.StatusCode, resp.Header, respBody)
}
