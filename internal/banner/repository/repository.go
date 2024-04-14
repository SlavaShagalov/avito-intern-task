package repository

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
)

type CreateParams struct {
	TagIDs    []int64
	FeatureID int64
	Content   map[string]any
	IsActive  bool
}

func (p *CreateParams) Validate() error {
	if p.TagIDs == nil {
		return pErrors.ErrBadTagIDsField
	}
	for _, tagID := range p.TagIDs {
		if tagID <= 0 {
			return pErrors.ErrBadTagIDsField
		}
	}
	if p.FeatureID <= 0 {
		return pErrors.ErrBadFeatureIDField
	}
	if p.Content == nil {
		return pErrors.ErrBadContentField
	}
	return nil
}

type FilterParams struct {
	FeatureID int64
	TagID     int64
	Limit     int
	Offset    int
}

type GetParams struct {
	FeatureID int64
	TagID     int64
	IsAdmin   bool
}

type PartialUpdateParams struct {
	ID        int64
	TagIDs    []int64
	FeatureID *int64
	Content   map[string]any
	IsActive  *bool
}

func (p *PartialUpdateParams) Validate() error {
	if p.ID <= 0 {
		return pErrors.ErrBadBannerIDParam
	}
	if p.TagIDs != nil {
		for _, tagID := range p.TagIDs {
			if tagID <= 0 {
				return pErrors.ErrBadTagIDsField
			}
		}
	}
	if p.FeatureID != nil && *p.FeatureID <= 0 {
		return pErrors.ErrBadFeatureIDField
	}
	return nil
}

type Repository interface {
	Create(ctx context.Context, params *CreateParams) (int64, error)
	List(ctx context.Context, params *FilterParams) ([]models.Banner, error)
	Get(ctx context.Context, params *GetParams) (*models.Banner, error)
	PartialUpdate(ctx context.Context, params *PartialUpdateParams) error
	Delete(ctx context.Context, id int64) error
}
