package cache

import "context"

type Value struct {
	Code int `json:"code"`
	Body any `json:"body"`
}

type Cache interface {
	Set(ctx context.Context, key string, value *Value) error
	Get(ctx context.Context, key string) (*Value, error)
}
