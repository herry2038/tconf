package backends

import "context"

// Response represents a response from a backend store.
type Response struct {
	Value []byte
	Error error
}

type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Watch(ctx context.Context, key string, stop chan bool) <-chan *Response
	Close() error
}
