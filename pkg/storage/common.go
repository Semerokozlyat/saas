package storage

type Storage interface {
	Get(id string) (interface{}, error)
	Put(interface{}) error
	Update(id string) error
	Delete(id string) error
}
