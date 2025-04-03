package services

import (
	"EfectiveMobile/internal/dto"
	"EfectiveMobile/internal/models"
	"EfectiveMobile/internal/repositories"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const (
	operatorIs   = "is"
	operatorIsnt = "isnt"
	operatorLs   = "ls"
	operatorMt   = "mt"

	apiGetAge         = "https://api.agify.io/?name=%s"
	apiGetGender      = "https://api.genderize.io/?name=%s"
	apiGetNationality = "https://api.nationalize.io/?name=%s"
)

type PersonService struct {
	PersonRepo *repositories.PersonRepo
	Log        *slog.Logger
}

func (ps *PersonService) GetPersonsByID(id int) (*models.Person, error) {
	p := models.Person{ID: id}
	return ps.PersonRepo.GetPersonByID(id, &p)
}

func (ps *PersonService) GetPersonsByParams(filters dto.Filters) ([]models.Person, error) {
	params := []string{}
	if filters.ByName != "" {
		validate := strings.Split(filters.ByName, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND name = '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'name is'", slog.String("name", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND name != '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'name is not'", slog.String("name", validate[1]))
		default:
			return nil, fmt.Errorf("invalid name param")
		}
	}
	if filters.BySurname != "" {
		validate := strings.Split(filters.BySurname, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND surname = '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'surname is'", slog.String("surname", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND surname != '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'surname is not'", slog.String("surname", validate[1]))
		default:
			return nil, fmt.Errorf("invalid surname param")
		}
	}
	if filters.ByPatronymic != "" {
		validate := strings.Split(filters.ByPatronymic, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND patronymic = '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'patronymic is'", slog.String("patronymic", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND patronymic != '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'patronymic is not'", slog.String("patronymic", validate[1]))
		default:
			return nil, fmt.Errorf("invalid patronymic param")
		}
	}
	if filters.ByAge != "" {
		validate := strings.Split(filters.ByAge, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND age = %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age is'", slog.String("age", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND age != %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age is not'", slog.String("age", validate[1]))
		case operatorLs:
			params = append(params, fmt.Sprintf("AND age < %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age less'", slog.String("age", validate[1]))
		case operatorMt:
			params = append(params, fmt.Sprintf("AND age > %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age more'", slog.String("age", validate[1]))
		default:
			return nil, fmt.Errorf("invalid age param")
		}
	}
	if filters.ByGender != "" {
		validate := strings.Split(filters.ByGender, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND gender = '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'gender is'", slog.String("gender", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND gender != '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'gender is not'", slog.String("gender", validate[1]))
		default:
			return nil, fmt.Errorf("invalid gender param")
		}
	}
	if filters.ByNationality != "" {
		validate := strings.Split(filters.ByNationality, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND nationality = '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'nationality is'", slog.String("nationality", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND nationality != '%s'", validate[1]))
			ps.Log.Debug("added filter parametr 'nationality is not'", slog.String("nationality", validate[1]))
		default:
			return nil, fmt.Errorf("invalid nationality param")
		}
	}
	if filters.ByLimit != 0 {
		params = append(params, fmt.Sprintf("LIMIT %d", filters.ByLimit))
		ps.Log.Debug("added filter parametr 'limit'", slog.Int("limit", filters.ByLimit))
	}
	if filters.ByOffset != 0 {
		params = append(params, fmt.Sprintf("OFFSET %d", filters.ByOffset))
		ps.Log.Debug("added filter parametr 'offset'", slog.Int("offset", filters.ByOffset))
	}
	filter := strings.Join(params, " ")

	return ps.PersonRepo.GetPersonsByParams(filter)
}

func (ps *PersonService) CreatePerson(person *dto.CreatePerson) (int, error) {
	for _, r := range person.Name {
		if !unicode.Is(unicode.Latin, r) {
			return 0, fmt.Errorf("name must be latin")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userData, err := getPersonData(ctx, person)
	if err != nil {
		return 0, err
	}
	ps.Log.Debug("get person data from api", slog.Any("person data", userData))

	return ps.PersonRepo.CreatePerson(userData)
}

func (ps *PersonService) DeletePersonById(id int) error {
	return ps.PersonRepo.DeletePersonById(id)
}

func (ps *PersonService) UpdatePerson(personDTO *dto.PersonUpdate) error {
	person, err := ps.GetPersonsByID(personDTO.ID)
	if err != nil {
		return err
	}
	ps.Log.Debug("get person to update", slog.Any("person", person))

	if personDTO.Name != "" {
		ps.Log.Debug("Name requires updated")
		person.Name = personDTO.Name
	}
	if personDTO.Surname != "" {
		ps.Log.Debug("Surname requires updated")
		person.Surname = personDTO.Surname
	}
	if personDTO.Patronymic != "" {
		ps.Log.Debug("Patronymic requires updated")
		person.Patronymic = personDTO.Patronymic
	}
	if personDTO.Age != 0 {
		ps.Log.Debug("Age requires updated")
		person.Age = personDTO.Age
	}
	if personDTO.Gender != "" {
		ps.Log.Debug("Gender requires updated")
		person.Gender = personDTO.Gender
	}
	if personDTO.Nationality != "" {
		ps.Log.Debug("Nationality requires updated")
		person.Nationality = personDTO.Nationality
	}

	return ps.PersonRepo.UpdatePerson(person)
}

func getPersonData(ctx context.Context, createdData *dto.CreatePerson) (*models.Person, error) {
	person := models.Person{Name: createdData.Name, Surname: createdData.Surname, Patronymic: createdData.Patronymic}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(apiGetAge, person.Name), nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		var data struct {
			Age int `json:"age"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		person.Age = data.Age
	} else {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout exceeded while getting age")
		}
		return nil, fmt.Errorf("cannot get age: %w", err)
	}

	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(apiGetGender, person.Name), nil)
	resp, err = http.DefaultClient.Do(req)
	if err == nil {
		var data struct {
			Gender string `json:"gender"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		person.Gender = data.Gender
	} else {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout exceeded while getting gender")
		}
		return nil, fmt.Errorf("cannot get gender: %w", err)
	}

	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(apiGetNationality, person.Name), nil)
	resp, err = http.DefaultClient.Do(req)
	if err == nil {
		var data struct {
			Country []struct {
				CountryID string `json:"country_id"`
			} `json:"country"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		if len(data.Country) > 0 {
			person.Nationality = data.Country[0].CountryID
		}
	} else {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout exceeded while getting 	nationality")
		}
		return nil, fmt.Errorf("cannot get nationality: %w", err)
	}
	return &person, nil
}
