package main

import (
	"EfectiveMobile/internal/db"
	"EfectiveMobile/internal/handlers"
	"EfectiveMobile/internal/repositories"
	"EfectiveMobile/internal/services"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// TODO Добавить адекватные информативные ошибки, а то щас какая-то тотал хуйня
// TODO Добавить логи
// TODO Добавить миграции

const (
	UsernameDB = "postgres"
	PasswordDB = "admin"
	HostDB     = "localhost"
	PortDB     = "5432"
	NameDB     = "efectivemobile"
)

func main() {
	conn, err := db.CreatePsqlConnection(db.ConnectionInfo{
		Username: UsernameDB,
		Password: PasswordDB,
		Host:     HostDB,
		Port:     PortDB,
		DBName:   NameDB,
	})

	if err != nil {
		log.Println("Failed connect to db")
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	router := chi.NewRouter()

	pr := &repositories.PersonRepo{DB: conn}
	ps := &services.PersonService{PersonRepo: pr}
	ph := handlers.PersonHandler{PersonService: ps}

	ph.Register(router)

	err = http.ListenAndServe(":8083", router)
	if err != nil {
		log.Println("Crashed server")
		os.Exit(1)
	}

}
