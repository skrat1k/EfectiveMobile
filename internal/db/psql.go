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

func CreatePsqlConnection(cfg ConnectionInfo) (*pgx.Conn, error) {
	connectionURL := MakeConnectionURL(cfg)

	conn, err := pgx.Connect(context.Background(), connectionURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func MakeConnectionURL(info ConnectionInfo) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", info.Username, info.Password, info.Host, info.Port, info.DBName)
}
