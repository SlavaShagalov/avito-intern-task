package bcrypt

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	pHasher "github.com/SlavaShagalov/avito-intern-task/internal/pkg/hasher"
)

type hasher struct{}

func New() pHasher.Hasher {
	return &hasher{}
}

func (h *hasher) GetHashedPassword(_ context.Context, password string) (string, error) {
	pswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pswd), err
}

func (h *hasher) CompareHashAndPassword(_ context.Context, hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
