package context

// The context package wraps the standard context package to add 
// user information, stored in private keys, to requests.

import (
	"context"

	"gastb.ar/models"
)

type privateKey string

// Declare unexported private keys
const (
	userKey privateKey = "user"
)

// WithUser adds user information to context.userKey
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User allows user information to be read from context
func User(ctx context.Context) *models.User {
	if check := ctx.Value(userKey); check != nil {
		if user, ok := check.(*models.User); ok {
			return user
		}
	}
	return nil
}
