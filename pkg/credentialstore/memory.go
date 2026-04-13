package credentialstore

import "fmt"

type memoryStore struct {
	items map[string][]byte
}

func NewMemory() Store {
	return &memoryStore{
		items: map[string][]byte{},
	}
}

func (m *memoryStore) Get(key string) ([]byte, error) {
	item, ok := m.items[key]
	if !ok {
		return nil, fmt.Errorf("failed to get item: item not found")
	}

	return append([]byte{}, item...), nil
}

func (m *memoryStore) Set(key string, data []byte) error {
	m.items[key] = append([]byte{}, data...)

	return nil
}

func (m *memoryStore) Remove(key string) error {
	delete(m.items, key)

	return nil
}
