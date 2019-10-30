package storage

type Storage interface {
	Get(id string) ([]byte, error)
	Put(string, []byte) error
	Update(id string) error
	Delete(id string) (int, error)
}
