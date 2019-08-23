package backends

// Response represents a response from a backend store.
type Response struct {
	Value []byte
	Error error
}

// KVPair holds both a key and value when reading a list.
type KVPair struct {
	Key   string
	Value []byte
}

type Store interface {
	Get(key string) ([]byte, error)
	Watch(key string, stop chan bool) <-chan *Response
}
