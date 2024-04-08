package models

import (
	"time"
)

type Banner struct {
	ID        int
	TagIDs    []int
	FeatureID int
	Content   map[string]any
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
