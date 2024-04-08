package auth

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
)

type SignInParams struct {
	Username string
	Password string
}

type SignUpParams struct {
	Username string
	Password string
}

type Usecase interface {
	SignIn(ctx context.Context, params *SignInParams) (*models.User, string, error)
	SignUp(ctx context.Context, params *SignUpParams) (*models.User, string, error)
}
