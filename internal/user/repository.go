package user

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
)

type CreateParams struct {
	Username string
	Password string
}

type Repository interface {
	Create(ctx context.Context, params *CreateParams) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}
