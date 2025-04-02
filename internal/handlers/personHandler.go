package handlers

import (
	"EfectiveMobile/internal/dto"
	"EfectiveMobile/internal/models"
	"EfectiveMobile/internal/services"
	"encoding/json"
	"log/slog"
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
	Log           *slog.Logger
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
		ph.Log.Error("Cannot get id", slog.String("error", err.Error()))
		return
	}
	ph.Log.Debug("Getting id", slog.Int("id", id))

	person, err := ph.PersonService.GetPersonsByID(id)
	if err != nil {
		http.Error(w, "Failed to get person", http.StatusInternalServerError)
		ph.Log.Error("Cannot get person by id", slog.Int("id", id), slog.String("error", err.Error()))
		return
	}
	json.NewEncoder(w).Encode(person)
	ph.Log.Debug("Encoded person to json", slog.Int("id", id))
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
			ph.Log.Error("Cannot get limit", slog.String("error", err.Error()))
			return
		}
		filters.ByLimit = limit
	}

	offsetStr := queryParams.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid offset value", http.StatusBadRequest)
			ph.Log.Error("Cannot get offset", slog.String("error", err.Error()))
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
		ph.Log.Error("Cannot decoded person to json", slog.String("error", err.Error()))
		return
	}

	id, err := ph.PersonService.CreatePerson(&person)
	if err != nil {
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		ph.Log.Error("Failed to create person", slog.String("error", err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, "Falider to encode JSON", http.StatusInternalServerError)
		ph.Log.Error("Failed to encode JSON", slog.String("error", err.Error()))
		return
	}
}

func (ph *PersonHandler) DeletePersonById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		ph.Log.Error("Cannot get id", slog.String("error", err.Error()))
		return
	}

	ph.Log.Debug("Getting id", slog.Int("id", id))

	err = ph.PersonService.DeletePersonById(id)
	if err != nil {
		http.Error(w, "Failed to delete person", http.StatusInternalServerError)
		ph.Log.Error("Failed to delete person", slog.String("error", err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ph *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	newData := dto.PersonUpdate{}
	err := json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		ph.Log.Error("Cannot decoded person to json", slog.String("error", err.Error()))
	}
	err = ph.PersonService.UpdatePerson(&newData)
	if err != nil {
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		ph.Log.Error("Failed to update person", slog.String("error", err.Error()))
	}

	w.WriteHeader(http.StatusNoContent)
}
