package storage

import (
	"fmt"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/jmoiron/sqlx"
	"log"
)

func NewPostgresStorage(cfg config.ServiceConfiguration) *PostgresStorage {

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		log.Fatalln(err)
	}

	db.DB.SetMaxOpenConns(40)
	db.DB.SetMaxIdleConns(10)
	db.DB.SetConnMaxLifetime(0)

	return &PostgresStorage{
		db: db,
	}
}

type PostgresStorage struct {
	db *sqlx.DB
}
