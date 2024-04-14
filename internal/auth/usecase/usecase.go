package usecase

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/auth"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/config"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	pHasher "github.com/SlavaShagalov/avito-intern-task/internal/pkg/hasher"
	bcryptHasher "github.com/SlavaShagalov/avito-intern-task/internal/pkg/hasher/bcrypt"
	"github.com/SlavaShagalov/avito-intern-task/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

const tokenExpiration = 30 * 24 * time.Hour

type usecase struct {
	usersRepo user.Repository
	hasher    pHasher.Hasher
	log       *zap.Logger
}

func New(usersRepo user.Repository, log *zap.Logger) auth.Usecase {
	return &usecase{
		usersRepo: usersRepo,
		hasher:    bcryptHasher.New(),
		log:       log,
	}
}

func (uc *usecase) SignIn(ctx context.Context, params *auth.SignInParams) (*models.User, string, error) {
	user, err := uc.usersRepo.GetByUsername(ctx, params.Username)
	if err != nil {
		return nil, "", err
	}

	if err = uc.hasher.CompareHashAndPassword(ctx, user.Password, params.Password); err != nil {
		return nil, "", errors.Wrap(pErrors.ErrWrongLoginOrPassword, err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(tokenExpiration).Unix(),
	})

	signedString, err := token.SignedString([]byte(viper.GetString(config.AuthKey)))
	if err != nil {
		return nil, "", err
	}

	uc.log.Debug("Sign In", zap.Int64("user_id", user.ID))
	return user, signedString, nil
}

func (uc *usecase) SignUp(ctx context.Context, params *auth.SignUpParams) (*models.User, string, error) {
	_, err := uc.usersRepo.GetByUsername(ctx, params.Username)
	if !errors.Is(err, pErrors.ErrUserNotFound) {
		if err != nil {
			return nil, "", err
		}
		return nil, "", pErrors.ErrUserAlreadyExists
	}

	hashedPassword, err := uc.hasher.GetHashedPassword(ctx, params.Password)
	if err != nil {
		return nil, "", pErrors.ErrGetHashedPassword
	}

	repParams := user.CreateParams{
		Username: params.Username,
		Password: hashedPassword,
	}
	user, err := uc.usersRepo.Create(ctx, &repParams)
	if err != nil {
		return nil, "", err
	}

	payload := jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedString, err := token.SignedString([]byte(viper.GetString(config.AuthKey)))
	if err != nil {
		return nil, "", err
	}

	uc.log.Debug("Sign Up", zap.Int64("user_id", user.ID))
	return user, signedString, nil
}
