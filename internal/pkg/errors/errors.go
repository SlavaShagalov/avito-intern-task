package errors

import (
	"errors"
)

var (
	ErrDb = errors.New("database error")

	// Banner
	ErrBannerNotFound = errors.New("banner not found")

	// User
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Auth
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrGetHashedPassword    = errors.New("get hashed password error")
	ErrInvalidAuthToken     = errors.New("invalid auth token")
	ErrAuthTokenNotFound    = errors.New("auth token not found")
	ErrAdminRequired        = errors.New("permission denied, this action requires admin privileges")

	// HTTP
	ErrReadBody         = errors.New("read request body error")
	ErrBadSessionCookie = errors.New("bad session cookie")

	ErrBadFeatureIDParam = errors.New("bad feature id parameter")
	ErrBadTagIDParam     = errors.New("bad tag id parameter")
	ErrBadLimitParam     = errors.New("bad limit parameter")
	ErrBadOffsetParam    = errors.New("bad offset parameter")
)
