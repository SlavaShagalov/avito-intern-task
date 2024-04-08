package boards

import (
	"context"
	pBannerRepo "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
)

type Usecase interface {
	Create(ctx context.Context, params *pBannerRepo.CreateParams) (int64, error)
	List(ctx context.Context, params *pBannerRepo.FilterParams) ([]models.Banner, error)
	Get(ctx context.Context, params *pBannerRepo.GetParams) (map[string]any, error)
	PartialUpdate(ctx context.Context, params *pBannerRepo.PartialUpdateParams) error
	Delete(ctx context.Context, id int64) error
}
