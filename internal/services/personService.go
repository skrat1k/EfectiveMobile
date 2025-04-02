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
			ps.Log.Debug("added filter parametr 'name is'")
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND name != %s", validate[1]))
			ps.Log.Debug("added filter parametr 'name is not'")
		default:
			return nil, fmt.Errorf("invalid name param")
		}
	}
	if filters.BySurname != "" {
		validate := strings.Split(filters.BySurname, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND surname = %s", validate[1]))
			ps.Log.Debug("added filter parametr 'surname is'")
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND surname != %s", validate[1]))
			ps.Log.Debug("added filter parametr 'surname is not'")
		default:
			return nil, fmt.Errorf("invalid surname param")
		}
	}
	if filters.ByPatronymic != "" {
		validate := strings.Split(filters.ByPatronymic, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND patronymic = %s", validate[1]))
			ps.Log.Debug("added filter parametr 'patronymic is'")
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND patronymic != %s", validate[1]))
			ps.Log.Debug("added filter parametr 'patronymic is not'")
		default:
			return nil, fmt.Errorf("invalid patronymic param")
		}
	}
	if filters.ByAge != "" {
		validate := strings.Split(filters.ByAge, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND age = %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age is'")
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND age != %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age is not'")
		case operatorLs:
			params = append(params, fmt.Sprintf("AND age < %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age less'")
		case operatorMt:
			params = append(params, fmt.Sprintf("AND age > %s", validate[1]))
			ps.Log.Debug("added filter parametr 'age more'")
		default:
			return nil, fmt.Errorf("invalid age param")
		}
	}
	if filters.ByGender != "" {
		validate := strings.Split(filters.ByGender, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND gender = %s", validate[1]))
			ps.Log.Debug("added filter parametr 'gender is'")
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND gender != %s", validate[1]))
			ps.Log.Debug("added filter parametr 'gender is not'")
		default:
			return nil, fmt.Errorf("invalid gender param")
		}
	}
	if filters.ByNationality != "" {
		validate := strings.Split(filters.ByNationality, ":")
		switch validate[0] {
		case operatorIs:
			params = append(params, fmt.Sprintf("AND nationality = %s", validate[1]))
			ps.Log.Debug("added filter parametr 'nationality is'")
		case operatorIsnt:
			params = append(params, fmt.Sprintf("AND nationality != %s", validate[1]))
			ps.Log.Debug("added filter parametr 'nationality is not'")
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
	userData, err := getPersonData(person)
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

func getPersonData(createdData *dto.CreatePerson) (*models.Person, error) {

	// TODO Добавить контекст с дедлайном, чтобы если какой-то из сервисов упадёт и я не смогу к нему достучаться - не зависнуть навечно
	person := models.Person{Name: createdData.Name, Surname: createdData.Surname, Patronymic: createdData.Patronymic}
	resp, err := http.Get(fmt.Sprintf(apiGetAge, person.Name))

	if err == nil {
		var data struct {
			Age int `json:"age"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		person.Age = data.Age
	} else {
		return nil, fmt.Errorf("сannot get age")
	}

	resp, err = http.Get(fmt.Sprintf(apiGetGender, person.Name))
	if err == nil {
		var data struct {
			Gender string `json:"gender"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		person.Gender = data.Gender
	} else {
		return nil, fmt.Errorf("сannot get gender")
	}

	resp, err = http.Get(fmt.Sprintf(apiGetNationality, person.Name))
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
		return nil, fmt.Errorf("сannot get nationality")
	}
	return &person, nil
}
