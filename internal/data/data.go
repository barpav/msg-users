package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db  *sql.DB
	cfg *Config
}

func (s *Storage) Init() (err error) {
	s.cfg = &Config{}
	s.cfg.Read()
	return s.connectToDatabase()
}

func (s *Storage) connectToDatabase() (err error) {
	dbAdress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", s.cfg.user, s.cfg.password, s.cfg.host, s.cfg.port, s.cfg.database)

	s.db, err = sql.Open("pgx", dbAdress)

	if err == nil {
		err = s.db.Ping()
	}

	if err == nil {
		log.Printf("Successfully connected to DB at %s", dbAdress)
	}

	return err
}
