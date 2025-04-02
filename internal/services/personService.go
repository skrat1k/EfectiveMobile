package services

import (
	"EfectiveMobile/internal/dto"
	"EfectiveMobile/internal/models"
	"EfectiveMobile/internal/repositories"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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

type userData struct {
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"country_id"`
}

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
			params = append(params, fmt.Sprintf("AND name = %s", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND name != %s", validate[1]))
		default:
			return nil, fmt.Errorf("invalid name param")
		}
	}
	if filters.BySurname != "" {
		validate := strings.Split(filters.BySurname, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND surname = %s", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND surname != %s", validate[1]))
		default:
			return nil, fmt.Errorf("invalid surname param")
		}
	}
	if filters.ByPatronymic != "" {
		validate := strings.Split(filters.ByPatronymic, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND patronymic = %s", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND patronymic != %s", validate[1]))
		default:
			return nil, fmt.Errorf("invalid patronymic param")
		}
	}
	if filters.ByAge != "" {
		validate := strings.Split(filters.ByAge, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND age = %s", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND age != %s", validate[1]))
		case operatorLs:
			params = append(params, fmt.Sprintf("AND age < %s", validate[1]))
		case operatorMt:
			params = append(params, fmt.Sprintf("AND age > %s", validate[1]))
		default:
			return nil, fmt.Errorf("invalid age param")
		}
	}
	if filters.ByGender != "" {
		validate := strings.Split(filters.ByGender, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND gender = %s", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND gender != %s", validate[1]))
		default:
			return nil, fmt.Errorf("invalid gender param")
		}
	}
	if filters.ByNationality != "" {
		validate := strings.Split(filters.ByNationality, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND nationality = %s", validate[1]))
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND nationality != %s", validate[1]))
		default:
			return nil, fmt.Errorf("invalid nationality param")
		}
	}
	if filters.ByLimit != 0 {
		params = append(params, fmt.Sprintf("LIMIT %d", filters.ByLimit))
	}
	if filters.ByOffset != 0 {
		params = append(params, fmt.Sprintf("OFFSET %d", filters.ByOffset))
	}
	filter := strings.Join(params, " ")

	return ps.PersonRepo.GetPersonsByParams(filter)
}

func (ps *PersonService) CreatePerson(person *models.Person) (int, error) {
	for _, r := range person.Name {
		if !unicode.Is(unicode.Latin, r) {
			return 0, fmt.Errorf("name must be latin")
		}
	}
	userData, err := getPersonData(person.Name)
	if err != nil {
		return 0, err
	}
	person.Age = userData.Age
	person.Gender = userData.Gender
	person.Nationality = userData.Nationality

	return ps.PersonRepo.CreatePerson(person)
}

func (ps *PersonService) DeletePersonById(id int) error {
	return ps.PersonRepo.DeletePersonById(id)
}

func (ps *PersonService) UpdatePerson(personDTO *dto.PersonUpdate) error {
	person, err := ps.GetPersonsByID(personDTO.ID)
	if err != nil {
		return err
	}

	if personDTO.Name != "" {
		person.Name = personDTO.Name
	}
	if personDTO.Surname != "" {
		person.Surname = personDTO.Surname
	}
	if personDTO.Patronymic != "" {
		person.Patronymic = personDTO.Patronymic
	}
	if personDTO.Age != 0 {
		person.Age = personDTO.Age
	}
	if personDTO.Gender != "" {
		person.Gender = personDTO.Gender
	}
	if personDTO.Nationality != "" {
		person.Nationality = personDTO.Nationality
	}

	return ps.PersonRepo.UpdatePerson(person)
}

func getPersonData(name string) (*userData, error) {

	// TODO Добавить контекст с дедлайном, чтобы если какой-то из сервисов упадёт и я не смогу к нему достучаться - не зависнуть навечно
	userData := userData{}
	resp, err := http.Get(fmt.Sprintf(apiGetAge, name))

	if err == nil {
		var data struct {
			Age int `json:"age"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		userData.Age = data.Age
	} else {
		return nil, fmt.Errorf("сannot get age")
	}

	resp, err = http.Get(fmt.Sprintf(apiGetGender, name))
	if err == nil {
		var data struct {
			Gender string `json:"gender"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		userData.Gender = data.Gender
	} else {
		return nil, fmt.Errorf("сannot get gender")
	}

	resp, err = http.Get(fmt.Sprintf(apiGetNationality, name))
	if err == nil {
		var data struct {
			Country []struct {
				CountryID string `json:"country_id"`
			} `json:"country"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		if len(data.Country) > 0 {
			userData.Nationality = data.Country[0].CountryID
		}
	} else {
		return nil, fmt.Errorf("сannot get nationality")
	}
	return &userData, nil
}
