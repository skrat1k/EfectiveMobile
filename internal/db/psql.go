package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ConnectionInfo struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

func CreatePsqlConnection(cfg string) (*pgx.Conn, error) {

	conn, err := pgx.Connect(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func MakeConnectionURL(info ConnectionInfo) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", info.Username, info.Password, info.Host, info.Port, info.DBName)
}
