package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db      *sql.DB
	cfg     *Config
	queries map[query]*sql.Stmt
}

type query interface {
	text() string
}

func queriesToPrepare() []query {
	return []query{
		queryValidateCredentials{},
		queryCreateUser{},
		queryGetUserInfoV1{},
		queryUpdateCommonProfileInfoV1{},
		queryChangePassword{},
	}
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
		for _, stmt := range s.queries {
			if stmt != nil {
				err = errors.Join(err, stmt.Close())
			}
		}

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
		log.Info().Msg(fmt.Sprintf("Successfully connected to DB at %s", dbAddress))
	}

	return err
}

func (s *Storage) prepareQueries() (err error) {
	s.queries = make(map[query]*sql.Stmt)

	for _, q := range queriesToPrepare() {
		err = errors.Join(err, s.prepare(q))
	}

	return err
}

func (s *Storage) prepare(q query) (err error) {
	s.queries[q], err = s.db.Prepare(q.text())
	return err
}
