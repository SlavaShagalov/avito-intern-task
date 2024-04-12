package http

import (
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	"time"
)

// API requests
type createRequest struct {
	TagIDs    []int64        `json:"tag_ids"`
	FeatureID int64          `json:"feature_id"`
	Content   map[string]any `json:"content"`
	IsActive  bool           `json:"is_active"`
}

type partialUpdateRequest struct {
	TagIDs    []int64        `json:"tag_ids"`
	FeatureID *int64         `json:"feature_id"`
	Content   map[string]any `json:"content"`
	IsActive  *bool          `json:"is_active"`
}

// API responses
type createResponse struct {
	BannerID int64 `json:"banner_id"`
}

func newCreateResponse(bannerID int64) *createResponse {
	return &createResponse{
		BannerID: bannerID,
	}
}

type banner struct {
	ID        int64          `json:"banner_id"`
	TagIDs    []int64        `json:"tag_ids"`
	FeatureID int64          `json:"feature_id"`
	Content   map[string]any `json:"content"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func newListResponse(banners []models.Banner) []banner {
	response := []banner{}
	for _, b := range banners {
		response = append(response, banner{
			ID:        b.ID,
			FeatureID: b.FeatureID,
			TagIDs:    b.TagIDs,
			Content:   b.Content,
			IsActive:  b.IsActive,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}
	return response
}
