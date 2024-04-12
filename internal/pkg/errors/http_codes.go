package errors

import "net/http"

var httpCodes = map[error]int{
	ErrDb: http.StatusInternalServerError,

	// Banner
	ErrBannerNotFound:      http.StatusNotFound,
	ErrBannerAlreadyExists: http.StatusBadRequest,

	// JSON
	ErrBadFeatureIDField: http.StatusBadRequest,
	ErrBadTagIDsField:    http.StatusBadRequest,
	ErrBadContentField:   http.StatusBadRequest,

	// User
	ErrUserNotFound:      http.StatusNotFound,
	ErrUserAlreadyExists: http.StatusConflict,

	// Auth
	ErrWrongLoginOrPassword: http.StatusBadRequest,
	ErrAuthTokenNotFound:    http.StatusUnauthorized,
	ErrInvalidAuthToken:     http.StatusUnauthorized,
	ErrAdminRequired:        http.StatusForbidden,

	// HTTP
	ErrReadBody: http.StatusBadRequest,

	// Get params
	ErrBadBannerIDParam:  http.StatusBadRequest,
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
