package handlers

import (
	"EfectiveMobile/internal/dto"
	"EfectiveMobile/internal/services"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	_ "EfectiveMobile/docs" // Подключаем документацию

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	httpSwagger "github.com/swaggo/http-swagger"
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
	ph.Log.Info("Successfully created http route", slog.String("route", getPersonByID))
	router.Get(getPersonByParams, ph.GetPersonsByParams)
	ph.Log.Info("Successfully created http route", slog.String("route", getPersonByParams))
	router.Delete(deletePersonByID, ph.DeletePersonById)
	ph.Log.Info("Successfully created http route", slog.String("route", deletePersonByID))
	router.Put(updatePerson, ph.UpdatePerson)
	ph.Log.Info("Successfully created http route", slog.String("route", updatePerson))
	router.Post(createPerson, ph.CreatePerson)
	ph.Log.Info("Successfully created http route", slog.String("route", createPerson))
	router.Get("/swagger/*", httpSwagger.WrapHandler)
	ph.Log.Info("Swagger documentation is enabled")
}

// @Summary Получение информации о человеке по ID
// @Description Возвращает данные о человеке по его идентификатору
// @Tags person
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} models.Person
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Failed to get person"
// @Router /api/v1/person/get/{id} [get]
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
		http.Error(w, fmt.Sprintf("Failed to get person: %s", err.Error()), http.StatusInternalServerError)
		ph.Log.Error("Cannot get person by id", slog.Int("id", id), slog.String("error", err.Error()))
		return
	}
	json.NewEncoder(w).Encode(person)
	ph.Log.Debug("Encoded person to json", slog.Int("id", id))
}

// @Summary Получение отфильтрованной информации о людях
// @Description Возвращает отфильтрованные данные о людях
// @Description Операторы для фильтрации значений (не распространяется на limit и offset):
// @Description - `var=is:X` — значение равно X
// @Description - `var=isnt:X` — значение не равно X
// @Description - `var=ls:X` — значение меньше X (только для age)
// @Description - `var=mt:X` — значение больше X (только для age)
// @Description - Пример:
// @Description - `age=mt:X` — значение больше X
// @Description - `name=is:X` — значение равно X
// @Tags person
// @Produce json
// @Param name query string false "Имя пользователя"
// @Param surname query string false "Фамилия пользователя"
// @Param patronymic query string false "Отчество пользователя"
// @Param gender query string false "Пол пользователя"
// @Param nationality query string false "Национальность пользователя"
// @Param age query int false "Возраст пользователя"
// @Param limit query int false "Лимит записей (если не задан - выводятся все подходящие данные)"
// @Param offset query int false "Смещение записей"
// @Success 200 {array} models.Person
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 500 {string} string "Failed to get persons"
// @Router /api/v1/person/get [get]
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
		http.Error(w, fmt.Sprintf("Failed to get person: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(persons)
}

// @Summary Создание нового пользователя
// @Description Создает нового пользователя с переданными данными
// @Tags person
// @Accept json
// @Produce json
// @Param person body dto.CreatePerson true "Данные пользователя для создания"
// @Success 201 {object} int "ID нового пользователя"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to create person"
// @Router /api/v1/person/create [post]
func (ph *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person dto.CreatePerson

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		ph.Log.Error("Cannot decoded person to json", slog.String("error", err.Error()))
		return
	}

	validate := validator.New()
	err = validate.Struct(person)
	if err != nil {
		http.Error(w, "Validation error: name and surname are required", http.StatusBadRequest)
		ph.Log.Error("Validation error", slog.String("error", err.Error()))
		return
	}

	id, err := ph.PersonService.CreatePerson(&person)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create person: %s", err.Error()), http.StatusInternalServerError)
		ph.Log.Error("Failed to create person", slog.String("error", err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
		ph.Log.Error("Failed to encode JSON", slog.String("error", err.Error()))
		return
	}
}

// @Summary Удаление пользователя по ID
// @Description Удаляет пользователя по переданному ID
// @Tags person
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 204 {string} string "User successfully deleted"
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Failed to delete person"
// @Router /api/v1/person/delete/{id} [delete]
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
		http.Error(w, fmt.Sprintf("Failed to delete person: %s", err.Error()), http.StatusInternalServerError)
		ph.Log.Error("Failed to delete person", slog.String("error", err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Обновление данных пользователя
// @Description Обновляет данные пользователя с переданными новыми данными
// @Tags person
// @Accept json
// @Produce json
// @Param person body dto.PersonUpdate true "Новые данные пользователя"
// @Success 204 {string} string "User successfully updated"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to update person"
// @Router /api/v1/person/update [put]
func (ph *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	newData := dto.PersonUpdate{}
	err := json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		ph.Log.Error("Cannot decoded person to json", slog.String("error", err.Error()))
	}
	err = ph.PersonService.UpdatePerson(&newData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update person: %s", err.Error()), http.StatusInternalServerError)
		ph.Log.Error("Failed to update person", slog.String("error", err.Error()))
	}

	w.WriteHeader(http.StatusNoContent)
}
