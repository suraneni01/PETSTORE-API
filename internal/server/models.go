package server

import (
	"errors"
	"fmt"
	"pet-store-server/internal/petstore"
)

type NewPetRequest struct {
	Type  petstore.PetType `json:"type"`
	Price float32          `json:"price"`
}

func (r *NewPetRequest) Validate() error {
	switch {
	case !r.Type.Valid():
		return fmt.Errorf("inappropriate pet type %s", r.Type)
	case r.Price < 0:
		return errors.New("price cannot be below zero")
	}
	return nil
}
