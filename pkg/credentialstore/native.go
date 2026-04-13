package credentialstore

import (
	"fmt"

	"github.com/99designs/keyring"
)

type nativeStore struct {
	keyring keyring.Keyring
}

type NativeArgs struct {
	Keyring keyring.Keyring
}

func NewNative(args NativeArgs) Store {
	return &nativeStore{
		keyring: args.Keyring,
	}
}

func (s *nativeStore) Get(key string) ([]byte, error) {
	item, err := s.keyring.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item.Data, nil
}

func (s *nativeStore) Set(key string, data []byte) error {
	err := s.keyring.Set(keyring.Item{
		Key:  key,
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("failed to set item: %w", err)
	}

	return nil
}

func (s *nativeStore) Remove(key string) error {
	err := s.keyring.Remove(key)
	if err != nil {
		return fmt.Errorf("failed to remove item: %w", err)
	}

	return nil
}
