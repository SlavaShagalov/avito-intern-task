package repository

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
)

type CreateParams struct {
	TagIDs    []int
	FeatureID int
	Content   map[string]any
	IsActive  bool
}

type FilterParams struct {
	FeatureID int
	TagID     int
	Limit     int
	Offset    int
}

type GetParams struct {
	FeatureID int
	TagID     int
}

type PartialUpdateParams struct {
	ID        int64
	TagIDs    []int
	FeatureID *int
	Content   map[string]any
	IsActive  *bool
}

type Repository interface {
	Create(ctx context.Context, params *CreateParams) (int64, error)
	List(ctx context.Context, params *FilterParams) ([]models.Banner, error)
	Get(ctx context.Context, params *GetParams) (map[string]any, error)
	PartialUpdate(ctx context.Context, params *PartialUpdateParams) error
	Delete(ctx context.Context, id int64) error
}
