package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/ognick/word_of_wisdom/pkg/logger"
)

type Client struct {
	addr    string
	handler func(conn net.Conn) error
	dialer  net.Dialer
	log     logger.Logger
}

func NewClient(
	addr string,
	handler func(conn net.Conn) error,
) *Client {
	return &Client{
		addr:    addr,
		handler: handler,
		log:     logger.NewLogger(),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	conn, err := c.dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			c.log.Errorf("failed to close: %v", err)
		}
	}()

	if err := c.handler(conn); err != nil {
		c.log.Errorf("failed to handle: %v", err)
	}

	return nil
}
