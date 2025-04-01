package repositories

import (
	"EfectiveMobile/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type PersonRepo struct {
	DB *pgx.Conn
}

func (pr *PersonRepo) GetPersonByID(id int, p *models.Person) (*models.Person, error) {
	query := "SELECT name, surname, COALESCE(patronymic, ''), age, gender, nationality FROM persons WHERE personid = $1"

	err := pr.DB.QueryRow(context.Background(), query, id).Scan(&p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
	if err != nil {
		return nil, err
	}
	return p, err
}

func (pr *PersonRepo) GetPersonsByParams(filter string) ([]models.Person, error) {
	query := "SELECT personid, name, surname, COALESCE(patronymic, ''), age, gender, nationality FROM persons WHERE 1=1 "
	if len(filter) > 0 {
		query = query + filter
	}

	rows, err := pr.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := []models.Person{}
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}

	return persons, nil

}

func (pr *PersonRepo) CreatePerson(person *models.Person) (int, error) {
	query := "INSERT INTO persons (name, surname, patronymic, age, gender, nationality) VALUES($1,$2,NULLIF($3, ''),$4,$5,$6) returning personid"
	var id int
	err := pr.DB.QueryRow(context.Background(), query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (pr *PersonRepo) DeletePersonById(id int) error {
	query := "DELETE FROM persons WHERE personid = $1"
	_, err := pr.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PersonRepo) UpdatePerson(person *models.Person) error {
	query := "UPDATE persons SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 WHERE personId = $7"
	_, err := pr.DB.Exec(context.Background(), query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, person.ID)
	if err != nil {
		return err
	}
	return nil
}
