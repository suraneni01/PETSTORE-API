package petstore

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination client_mock.go -package petstore . Interface
type Interface interface {
	Add(t PetType, price float32) error
	Get(id uint64) (*Pet, error)
	List() ([]*Pet, error)
}
