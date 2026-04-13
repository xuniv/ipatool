//go:build js && wasm

package credentialstore

import "fmt"

type HostSecretFunc func(key string) ([]byte, bool)

type WebArgs struct {
	HostSecret HostSecretFunc
}

type webStore struct {
	memory     Store
	hostSecret HostSecretFunc
}

func NewWeb(args WebArgs) Store {
	return &webStore{
		memory:     NewMemory(),
		hostSecret: args.HostSecret,
	}
}

func (w *webStore) Get(key string) ([]byte, error) {
	if item, err := w.memory.Get(key); err == nil {
		return item, nil
	}

	if w.hostSecret != nil {
		secret, ok := w.hostSecret(key)
		if ok {
			return append([]byte{}, secret...), nil
		}
	}

	return nil, fmt.Errorf("failed to get item: item not found")
}

func (w *webStore) Set(key string, data []byte) error {
	return w.memory.Set(key, data)
}

func (w *webStore) Remove(key string) error {
	return w.memory.Remove(key)
}
