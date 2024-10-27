package postgres

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var postgresContainerImg = "postgres:16"

// InitPostgresTestContainer is a helper function that starts and prepare a Postgres container ready for testing.
// When the function returns the database is ready and a connection string and a closing function are returned.
func InitPostgresTestContainer() (string, func() error, error) {
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(postgresContainerImg),
		postgres.WithUsername("testcontainer_user"),
		postgres.WithPassword("testcontainer_p@ssw0rd"),
		postgres.WithDatabase("myservice_users"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return "", nil, err
	}
	deferFunc := func() error {
		return pgContainer.Terminate(ctx)
	}
	conStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	return conStr, deferFunc, err
}
