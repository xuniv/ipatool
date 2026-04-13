package credentialstore

//go:generate go run go.uber.org/mock/mockgen -source=store.go -destination=store_mock.go -package credentialstore
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Remove(key string) error
}
