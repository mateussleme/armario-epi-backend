package database

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
)

var Db *sql.DB

func Init() error {
	port, err := strconv.Atoi(os.Getenv("SQL_PORT"))
	if err != nil {
		return err
	}

	cfg := pq.Config{
		Host:           os.Getenv("SQL_HOST"),
		Port:           uint16(port),
		User:           os.Getenv("SQL_USER"),
		Password:       os.Getenv("SQL_PASSWORD"),
		Database:       os.Getenv("SQL_DB"),
		ConnectTimeout: 5 * time.Second,
		SSLMode:        pq.SSLModePrefer,
	}

	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return err
	}

	db := sql.OpenDB(c)

	err = db.Ping()
	if err != nil {
		return err
	}

	Db = db

	return nil
}

func Connection(ctx context.Context) (*sql.Conn, error) {
	return Db.Conn(ctx)
}
