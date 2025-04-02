package repositories

import (
	"EfectiveMobile/internal/models"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type PersonRepo struct {
	DB  *pgx.Conn
	Log *slog.Logger
}

func (pr *PersonRepo) GetPersonByID(id int, p *models.Person) (*models.Person, error) {
	query := "SELECT name, surname, COALESCE(patronymic, ''), age, gender, nationality FROM persons WHERE personid = $1"
	pr.Log.Debug("Query to DB", slog.String("Query", query), slog.Int("personid", id))

	err := pr.DB.QueryRow(context.Background(), query, id).Scan(&p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
	if err != nil {
		return nil, err
	}
	pr.Log.Debug("Returning person", slog.Any("person", p))
	return p, err
}

func (pr *PersonRepo) GetPersonsByParams(filter string) ([]models.Person, error) {
	query := "SELECT personid, name, surname, COALESCE(patronymic, ''), age, gender, nationality FROM persons WHERE 1=1 "
	if len(filter) > 0 {
		query = query + filter
	}
	pr.Log.Debug("Query to DB with filter", slog.String("Query", query), slog.String("filter", filter))

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
		pr.Log.Debug("Add person to returning", slog.Any("person", p))
		persons = append(persons, p)
	}

	return persons, nil

}

func (pr *PersonRepo) CreatePerson(person *models.Person) (int, error) {
	query := "INSERT INTO persons (name, surname, patronymic, age, gender, nationality) VALUES($1,$2,NULLIF($3, ''),$4,$5,$6) returning personid"
	pr.Log.Debug("Query to create person", slog.String("Query", query))
	var id int
	err := pr.DB.QueryRow(context.Background(), query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&id)
	if err != nil {
		return 0, err
	}
	pr.Log.Debug("Succesful created person", slog.Any("person data", person))
	return id, nil
}

func (pr *PersonRepo) DeletePersonById(id int) error {
	query := "DELETE FROM persons WHERE personid = $1"
	pr.Log.Debug("Query to delete person", slog.String("Query", query))
	_, err := pr.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	pr.Log.Debug("Succesful delete person", slog.Int("personID", id))
	return nil
}

func (pr *PersonRepo) UpdatePerson(person *models.Person) error {
	query := "UPDATE persons SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 WHERE personId = $7"
	pr.Log.Debug("Query to delete person", slog.String("Query", query))
	_, err := pr.DB.Exec(context.Background(), query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, person.ID)
	if err != nil {
		return err
	}
	pr.Log.Debug("Succesful update person", slog.Any("person data", person))
	return nil
}
