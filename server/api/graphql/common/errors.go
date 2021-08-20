package common

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrAuthorizeRequired = errors.New("authorize required")
	ErrInvalidUUID       = errors.New("invalid UUID")
)
