package usecase

import (
	"context"
	pBanner "github.com/SlavaShagalov/avito-intern-task/internal/banner"
	"github.com/SlavaShagalov/avito-intern-task/internal/banner/repository"
	pBannerRepo "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	"go.uber.org/zap"
)

type usecase struct {
	repo repository.Repository
	log  *zap.Logger
}

func New(repo repository.Repository, log *zap.Logger) pBanner.Usecase {
	return &usecase{
		repo: repo,
		log:  log,
	}
}

func (uc *usecase) Create(ctx context.Context, params *pBannerRepo.CreateParams) (int64, error) {
	return uc.repo.Create(ctx, params)
}

func (uc *usecase) List(ctx context.Context, params *pBannerRepo.FilterParams) ([]models.Banner, error) {
	return uc.repo.List(ctx, params)
}

func (uc *usecase) Get(ctx context.Context, params *pBannerRepo.GetParams) (map[string]any, error) {
	banner, err := uc.repo.Get(ctx, params)
	if err != nil {
		return nil, err
	}
	if !banner.IsActive && !params.IsAdmin {
		return nil, pErrors.ErrBannerDisabled
	}
	return banner.Content, nil
}

func (uc *usecase) PartialUpdate(ctx context.Context, params *pBannerRepo.PartialUpdateParams) error {
	return uc.repo.PartialUpdate(ctx, params)
}

func (uc *usecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
