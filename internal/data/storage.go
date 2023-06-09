package data

import (
	"context"
	"database/sql"
	"errors"
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

func (s *Storage) Open() (err error) {
	s.cfg = &Config{}
	s.cfg.Read()

	err = s.connectToDatabase()

	if err != nil {
		return err
	}

	return s.prepareQueries()
}

func (s *Storage) Close(ctx context.Context) (err error) {
	closed := make(chan struct{}, 1)

	go func() {
		err = errors.Join(err, s.queries.register.Close())
		err = errors.Join(err, s.db.Close())

		closed <- struct{}{}
	}()

	select {
	case <-closed:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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
