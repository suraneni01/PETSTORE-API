package server

import (
	"encoding/json"
	"log"
	"net/http"
	"pet-store-server/internal/petstore"
	"strconv"

	"github.com/gorilla/mux"
)

type Service struct {
	petstore petstore.Interface
	router   *mux.Router
}

func NewService(petstore petstore.Interface) *Service {
	s := Service{
		petstore: petstore,
		router:   mux.NewRouter(),
	}

	s.router.Methods(http.MethodGet).Path("/pets").HandlerFunc(s.listPets)
	s.router.Methods(http.MethodGet).Path("/pets/{id:[0-9]+}").HandlerFunc(s.getPet)
	s.router.Methods(http.MethodPost).Path("/pets").HandlerFunc(s.addPet)

	return &s
}

func (s *Service) Start(addr string) error {
	log.Printf("starting server on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

func (s *Service) addPet(w http.ResponseWriter, r *http.Request) {
	var newPet NewPetRequest
	err := json.NewDecoder(r.Body).Decode(&newPet)
	if err != nil {
		s.serError(w, http.StatusInternalServerError, err)
		return
	}

	err = newPet.Validate()
	if err != nil {
		s.serError(w, http.StatusBadRequest, err)
		return
	}

	err = s.petstore.Add(newPet.Type, newPet.Price)
	if err != nil {
		s.serError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Service) getPet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		s.serError(w, http.StatusBadRequest, err)
		return
	}

	pet, err := s.petstore.Get(id)
	if err != nil {
		s.serError(w, http.StatusInternalServerError, err)
		return
	}

	data, err := json.Marshal(pet)
	if err != nil {
		s.serError(w, http.StatusInternalServerError, err)
		return
	}
	_, _ = w.Write(data)
}

func (s *Service) listPets(w http.ResponseWriter, r *http.Request) {
	pets, err := s.petstore.List()
	if err != nil {
		s.serError(w, http.StatusInternalServerError, err)
		return
	}

	data, err := json.Marshal(pets)
	if err != nil {
		s.serError(w, http.StatusInternalServerError, err)
		return
	}
	_, _ = w.Write(data)
}

func (s *Service) serError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(err.Error()))
}
