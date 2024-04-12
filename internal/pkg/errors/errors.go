package errors

import (
	"errors"
)

var (
	ErrDb      = errors.New("database error")
	ErrUnknown = errors.New("unknown error")

	// Banner
	ErrBannerNotFound      = errors.New("banner not found")
	ErrBannerAlreadyExists = errors.New("banner with such feature and tag already exists")

	// User
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Auth
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrGetHashedPassword    = errors.New("get hashed password error")
	ErrInvalidAuthToken     = errors.New("invalid auth token")
	ErrAuthTokenNotFound    = errors.New("auth token not found")

	// Access
	ErrAdminRequired  = errors.New("admin privileges required")
	ErrBannerDisabled = errors.New("banner disabled")

	// HTTP
	ErrReadBody = errors.New("read request body error")

	// JSON
	ErrBadContentField   = errors.New("bad content field")
	ErrBadFeatureIDField = errors.New("bad feature_id field")
	ErrBadTagIDsField    = errors.New("bad tag_ids field")

	// Get params
	ErrBadBannerIDParam  = errors.New("bad banner id parameter")
	ErrBadFeatureIDParam = errors.New("bad feature id parameter")
	ErrBadTagIDParam     = errors.New("bad tag id parameter")
	ErrBadLimitParam     = errors.New("bad limit parameter")
	ErrBadOffsetParam    = errors.New("bad offset parameter")
)
