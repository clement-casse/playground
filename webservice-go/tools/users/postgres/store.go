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
	"golang.org/x/crypto/bcrypt"

	"github.com/clement-casse/playground/webservice-go/tools/users"
)

var (
	bcryptPasswordCost = 8

	addUserStmt  = `INSERT INTO "users" (email, username, password) VALUES ($1, $2, $3)`
	getUserStmt  = `SELECT email, username FROM "users" WHERE email=$1`
	delUserStmt  = `DELETE FROM "users" WHERE email=$1`
	authUserStmt = `SELECT password FROM "users" WHERE email=$1`
)

type userStore struct {
	db *sql.DB
}

// NewUserStore created a User Store with a Postgres Database backend.
// All migrations are run by the store at creation time and are effective before returning.
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

func (s *userStore) CreateUser(ctx context.Context, name, email, password string) (*users.User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcryptPasswordCost)
	if err != nil {
		return nil, fmt.Errorf("postgresUserStore/CreateUser: %w", err)
	}

	result, err := s.db.ExecContext(ctx, addUserStmt, email, name, string(hashedPwd))
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
	if n, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("postgresUserStore/DeleteUser: %w", err)
	} else if n != 1 {
		return fmt.Errorf("postgresUserStore/DeleteUser: incorrect number of rows modified, expected 1, got %d", n)
	}
	return nil
}

// Authenticate performs a password authentication to the given user. Users are identified by their email.
// It expect only one param being the password.
func (s *userStore) Authenticate(ctx context.Context, userID string, params ...string) (*users.User, error) {
	if len(params) != 1 {
		return nil, errors.New("postgresUserStore/Authenticate: only expecting one param being the password")
	}
	credPassword := params[0]

	var userHashedPassword string

	if err := s.db.QueryRowContext(ctx, authUserStmt, userID).
		Scan(&userHashedPassword); errors.Is(err, sql.ErrNoRows) {
		return nil, users.ErrAuthenticationFailure // User simply does not exist
	} else if err != nil {
		return nil, fmt.Errorf("postgresUserStore/Authenticate: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userHashedPassword), []byte(credPassword)); err != nil {
		return nil, users.ErrAuthenticationFailure // Password does not match
	}

	var username, email string
	if err := s.db.QueryRowContext(ctx, getUserStmt, userID).
		Scan(&email, &username); err != nil {
		return nil, err
	}

	return &users.User{
		Name:  username,
		Email: email,
	}, nil
}
