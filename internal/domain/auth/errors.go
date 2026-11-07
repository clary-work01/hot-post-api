package auth

import "errors"

var (
	errInvalidRefreshToken = errors.New("auth: invalid refresh token")
	errInvalidCredentials  = errors.New("auth: invalid username or password")
)
