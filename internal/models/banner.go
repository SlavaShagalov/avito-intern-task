package models

import (
	"time"
)

type Banner struct {
	ID        int64
	TagIDs    []int64
	FeatureID int64
	Content   map[string]any
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
