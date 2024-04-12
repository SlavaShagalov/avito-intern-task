package errors

var jsonErrors = map[error]struct{}{
	ErrDb:      {},
	ErrUnknown: {},

	// Banner
	ErrBannerAlreadyExists: {},

	// JSON
	ErrBadFeatureIDField: {},
	ErrBadTagIDsField:    {},
	ErrBadContentField:   {},

	// User
	ErrUserAlreadyExists: {},

	// Auth
	ErrWrongLoginOrPassword: {},
	ErrAuthTokenNotFound:    {},

	// HTTP
	ErrReadBody: {},

	// Get params
	ErrBadBannerIDParam:  {},
	ErrBadFeatureIDParam: {},
	ErrBadTagIDParam:     {},
	ErrBadLimitParam:     {},
	ErrBadOffsetParam:    {},
}

func IsJSONError(err error) bool {
	_, exist := jsonErrors[err]
	return exist
}
