package main

import (
	"EfectiveMobile/internal/config"
	"EfectiveMobile/internal/db"
	"EfectiveMobile/internal/handlers"
	"EfectiveMobile/internal/repositories"
	"EfectiveMobile/internal/services"
	"EfectiveMobile/pkg/logger"
	"context"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// TODO Добавить адекватные информативные ошибки, а то щас какая-то тотал хуйня
// TODO Добавить логи
// TODO Потестить лимиты и оффсеты
// TODO Потестить код в целом
// TODO Добавить сваггер

func main() {

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err)
	}
	log := logger.GetLogger(cfg.Env)

	dbConnectionUrl := db.MakeConnectionURL(db.ConnectionInfo{
		Username: cfg.Username,
		Password: cfg.Password,
		Host:     cfg.Host,
		Port:     cfg.Port,
		DBName:   cfg.Name,
	})

	if err := db.RunMigrations(dbConnectionUrl); err != nil {
		if strings.Contains(err.Error(), "no change") {
			log.Info("No new migrations found, continuing application startup...")
		} else {
			log.Error("Migration error", slog.String("error", err.Error()))
			panic(err)
		}
	}

	conn, err := db.CreatePsqlConnection(dbConnectionUrl)

	if err != nil {
		log.Error("Failed connect to db", slog.String("error", err.Error()))
		panic(err)
	}

	defer conn.Close(context.Background())

	router := chi.NewRouter()

	pr := &repositories.PersonRepo{DB: conn}
	ps := &services.PersonService{PersonRepo: pr, Log: log}
	ph := handlers.PersonHandler{PersonService: ps, Log: log}

	ph.Register(router)

	err = http.ListenAndServe(":8083", router)
	if err != nil {
		log.Error("Crashed server", slog.String("error", err.Error()))
		panic(err)
	}

}
