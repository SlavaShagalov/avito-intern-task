package errors

import (
	"errors"
)

var (
	ErrDb = errors.New("database error")

	// Banner
	ErrBannerNotFound = errors.New("banner not found")

	// HTTP
	ErrReadBody         = errors.New("read request body error")
	ErrBadSessionCookie = errors.New("bad session cookie")

	ErrBadFeatureIDParam = errors.New("bad feature id parameter")
	ErrBadTagIDParam     = errors.New("bad tag id parameter")
	ErrBadLimitParam     = errors.New("bad limit parameter")
	ErrBadOffsetParam    = errors.New("bad offset parameter")
)
