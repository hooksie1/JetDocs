package backend

import "github.com/nats-io/nats.go"

const (
	natsBucket = "pages"
)

type Backend interface {
	Write(string, []byte) (int, error)
}

type Nats struct {
	Conn *nats.Conn
}

func (n Nats) Write(name string, data []byte) (int, error) {
	js, err := n.Conn.JetStream()
	if err != nil {
		return 0, err
	}

	kv, err := js.KeyValue(natsBucket)
	if err != nil {
		return 0, err
	}

	if _, err := kv.Put(name, data); err != nil {
		return 0, err
	}

	return len(data), nil
}
