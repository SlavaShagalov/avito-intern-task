package errors

import "net/http"

var httpCodes = map[error]int{
	ErrDb: http.StatusInternalServerError,

	// Banner
	ErrBannerNotFound: http.StatusNotFound,

	// HTTP
	ErrReadBody:         http.StatusBadRequest,
	ErrBadSessionCookie: http.StatusBadRequest,

	// Query params
	ErrBadFeatureIDParam: http.StatusBadRequest,
	ErrBadTagIDParam:     http.StatusBadRequest,
	ErrBadLimitParam:     http.StatusBadRequest,
	ErrBadOffsetParam:    http.StatusBadRequest,
}

func ErrorToHTTPCode(err error) (int, bool) {
	httpCode, exist := httpCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
