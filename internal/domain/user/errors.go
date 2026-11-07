package user

import "errors"

var (
	ErrUserNotFound     = errors.New("user: not found")
	errCannotFollowSelf = errors.New("user: cannot follow yourself")
	ErrUsernameTaken    = errors.New("user: user name taken")
)
