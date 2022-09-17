package counter

import (
	"context"
	"sync"
)

var _ counter = (*Client)(nil)

const (
	GC Name = iota
	ZNUM
	DC
	RCTVNUM
	RCTNUM
)

const (
	GET Operation = iota
	RESET
	INC
)

type (
	Config struct {
	}
	Client struct {
		conf *Config
		mu   *sync.Mutex
	}
	Name      int
	Operation int
	counter   interface {
		Do(context.Context, Name, Operation) (int64, error)
	}
)

func NewClient(cfg *Config) *Client {
	return &Client{
		mu:   &sync.Mutex{},
		conf: cfg,
	}
}

func (c *Client) Do(ctx context.Context, name Name, op Operation) (int64, error) {
	return 0, nil
}
