package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // loading file driver for migrations
	_ "github.com/lib/pq"                                // postgres driver for database/sql

	"github.com/clement-casse/playground/webservice-go/tools/users"
)

var (
	addUserStmt  = `INSERT INTO "users" (email, username) VALUES ($1, $2)`
	delUserStmt  = `DELETE FROM "users" ...`
	authUserStmt = `SELECT * FROM "users" WHERE ...`
)

type userStore struct {
	db *sql.DB
}

// NewUserStore created a User Store with a Postgres Database backend.
// migrations are run before returning.
func NewUserStore(conStr string) (users.Store, error) {
	var store = &userStore{}
	var err error
	store.db, err = sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}
	driver, err := pgmigrate.WithInstance(store.db, &pgmigrate.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}
	return store, m.Up()
}

func (s *userStore) CreateUser(ctx context.Context, name, email string) (*users.User, error) {
	result, err := s.db.ExecContext(ctx, addUserStmt, name, email)
	if err != nil {
		return nil, fmt.Errorf("postgresUserStore/CreateUser: %w", err)
	}

	if n, err := result.RowsAffected(); err != nil {
		return nil, fmt.Errorf("postgresUserStore/CreateUser: %w", err)
	} else if n != 1 {
		return nil, fmt.Errorf("postgresUserStore/CreateUser: incorrect number of rows modified, expected 1, got %d", n)
	}

	return &users.User{
		Name:  name,
		Email: email,
	}, nil
}

func (s *userStore) DeleteUser(ctx context.Context, user *users.User) error {
	result, err := s.db.ExecContext(ctx, delUserStmt, user.Email)
	if err != nil {
		return err
	}
	_ = result
	return errors.New("not implemented")
}

func (s *userStore) Authenticate(ctx context.Context, userID string, params ...string) (*users.User, error) {
	result, err := s.db.ExecContext(ctx, authUserStmt, userID)
	if err != nil {
		return nil, err
	}
	_ = result
	_ = params
	return nil, errors.New("not implemented")
}
