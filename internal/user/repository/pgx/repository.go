package pgx

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	pUsers "github.com/SlavaShagalov/avito-intern-task/internal/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func New(pool *pgxpool.Pool, log *zap.Logger) pUsers.Repository {
	return &repository{
		pool: pool,
		log:  log,
	}
}

const createCmd = `
	INSERT INTO users (username, password)
	VALUES ($1, $2)
	RETURNING id, username, password, created_at;`

func (repo *repository) Create(ctx context.Context, params *pUsers.CreateParams) (*models.User, error) {
	row := repo.pool.QueryRow(ctx, createCmd, params.Username, params.Password)

	user := new(models.User)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err))
		return nil, pErrors.ErrDb
	}

	repo.log.Debug("User created", zap.Int64("user_id", user.ID))
	return user, nil
}

const getByUsernameCmd = `
	SELECT id, username, password, is_admin, created_at
	FROM users
	WHERE username = $1;`

func (repo *repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	row := repo.pool.QueryRow(ctx, getByUsernameCmd, username)

	user := new(models.User)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pErrors.ErrUserNotFound
		}
		repo.log.Error(constants.DBError, zap.Error(err))
		return nil, pErrors.ErrDb
	}

	return user, nil
}
