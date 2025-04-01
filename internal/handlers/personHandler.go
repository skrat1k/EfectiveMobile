package handlers

import (
	"EfectiveMobile/internal/dto"
	"EfectiveMobile/internal/models"
	"EfectiveMobile/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	getPersonByID     = "/api/v1/person/get/{id}"
	getPersonByParams = "/api/v1/person/get"
	deletePersonByID  = "/api/v1/person/delete/{id}"
	updatePerson      = "/api/v1/person/update"
	createPerson      = "/api/v1/person/create"
)

type PersonHandler struct {
	PersonService *services.PersonService
}

func (ph *PersonHandler) Register(router *chi.Mux) {
	router.Get(getPersonByID, ph.GetPersonsByID)
	router.Get(getPersonByParams, ph.GetPersonsByParams)
	router.Delete(deletePersonByID, ph.DeletePersonById)
	router.Put(updatePerson, ph.UpdatePerson)
	router.Post(createPerson, ph.CreatePerson)
}

func (ph *PersonHandler) GetPersonsByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	person, err := ph.PersonService.GetPersonsByID(id)
	if err != nil {
		http.Error(w, "Failed to ", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (ph *PersonHandler) GetPersonsByParams(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filters := dto.Filters{}
	filters.ByName = queryParams.Get("name")
	filters.BySurname = queryParams.Get("surname")
	filters.ByPatronymic = queryParams.Get("patronymic")
	filters.ByGender = queryParams.Get("gender")
	filters.ByNationality = queryParams.Get("nationality")
	filters.ByAge = queryParams.Get("age")

	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit value", http.StatusBadRequest)
			return
		}
		filters.ByLimit = limit
	}

	offsetStr := queryParams.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid offset value", http.StatusBadRequest)
			return
		}
		filters.ByOffset = offset
	}

	persons, err := ph.PersonService.GetPersonsByParams(filters)
	if err != nil {
		http.Error(w, "Failed to get persons", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(persons)
}

func (ph *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := ph.PersonService.CreatePerson(&person)
	if err != nil {
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, "Falider to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (ph *PersonHandler) DeletePersonById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = ph.PersonService.DeletePersonById(id)
	if err != nil {
		http.Error(w, "Failed to delete person", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ph *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	newData := dto.PersonUpdate{}
	err := json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	err = ph.PersonService.UpdatePerson(&newData)
	if err != nil {
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
