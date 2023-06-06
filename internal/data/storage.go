package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db      *sql.DB
	cfg     *Config
	queries *queries
}

type queries struct {
	register *sql.Stmt
}

func (s *Storage) Init() (err error) {
	s.cfg = &Config{}
	s.cfg.Read()

	err = s.connectToDatabase()

	if err != nil {
		return err
	}

	return s.prepareQueries()
}

func (s *Storage) connectToDatabase() (err error) {
	dbAddress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", s.cfg.user, s.cfg.password, s.cfg.host, s.cfg.port, s.cfg.database)

	s.db, err = sql.Open("pgx", dbAddress)

	if err == nil {
		err = s.db.Ping()
	}

	if err == nil {
		log.Printf("Successfully connected to DB at %s", dbAddress)
	}

	return err
}

func (s *Storage) prepareQueries() (err error) {
	s.queries = &queries{}

	err = s.prepareRegisterQuery()

	return err
}
