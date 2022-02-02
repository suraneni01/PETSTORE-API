package petstore

type Pet struct {
	ID    uint64  `json:"id"`
	Type  PetType `json:"type"`
	Price float32 `json:"price"`
}

type PetType string

const (
	Dog   PetType = "dog"
	Cat   PetType = "cat"
	Fish  PetType = "fish"
	Bird  PetType = "bird"
	Gecko PetType = "gecko"
)

func (pt PetType) Valid() bool {
	switch pt {
	case Dog, Cat, Fish, Bird, Gecko:
		return true
	}
	return false
}
