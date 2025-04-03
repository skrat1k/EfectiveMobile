package main

import (
	"EfectiveMobile/internal/config"
	"EfectiveMobile/internal/db"
	"EfectiveMobile/internal/handlers"
	"EfectiveMobile/internal/repositories"
	"EfectiveMobile/internal/services"
	"EfectiveMobile/pkg/logger"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// @title EffectiveMobile API
// @version 1.0
// @description API для обогащений пользовательских данных
// @host localhost:8083
// @BasePath /
// @schemes http
func main() {

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err)
	}
	log := logger.GetLogger(cfg.Env)

	log.Info("Logger successfully initialized")

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
	log.Info("Successfully connect to db")

	router := chi.NewRouter()

	pr := &repositories.PersonRepo{DB: conn, Log: log}
	ps := &services.PersonService{PersonRepo: pr, Log: log}
	ph := handlers.PersonHandler{PersonService: ps, Log: log}

	ph.Register(router)

	log.Info("Starting server...", slog.String("Address", fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort), router)
	if err != nil {
		log.Error("Crashed server", slog.String("error", err.Error()))
		panic(err)
	}

}
