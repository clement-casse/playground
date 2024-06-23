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
	user1, err = userStore.CreateUser(ctx, "user1", "user1@domain.example", "user1Password")
	assert.NoError(t, err)
	assert.Equal(t, &users.User{Email: "user1@domain.example", Name: "user1"}, user1)

	notAUser, err := userStore.Authenticate(ctx, "DoesNotExist@domain.example", "doesItMatterAtThisPoint?")
	assert.ErrorIs(t, err, users.ErrAuthenticationFailure)
	assert.Nil(t, notAUser)

	notAUser, err = userStore.Authenticate(ctx, user1.Email, "wrongPassword")
	assert.ErrorIs(t, err, users.ErrAuthenticationFailure)
	assert.Nil(t, notAUser)

	anotherUser1, err := userStore.Authenticate(ctx, user1.Email, "user1Password")
	assert.NoError(t, err)
	assert.Equal(t, &users.User{Email: "user1@domain.example", Name: "user1"}, anotherUser1)
}
