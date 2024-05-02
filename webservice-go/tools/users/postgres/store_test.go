package postgres

import (
	"context"
	"testing"

	"github.com/clement-casse/playground/webservice-go/tools/users"
	"github.com/stretchr/testify/assert"
)

func TestUserStore(t *testing.T) {
	conStr, deferFunc, err := InitPostgresTestContainer()
	if err != nil {
		t.Fatalf("Test failed for unexpected reason %s", err)
	}
	defer func() {
		if deferErr := deferFunc(); deferErr != nil {
			t.Fatalf("unexpected err: %s", deferErr)
		}
	}()

	userStore, err := NewUserStore(conStr)
	assert.NoError(t, err)

	ctx := context.Background()
	var user1 *users.User
	user1, err = userStore.CreateUser(ctx, "user1", "user1@domain.example")
	assert.NoError(t, err)
	assert.Equal(t, &users.User{Email: "user1@domain.example", Name: "user1"}, user1)
}
