package user

import "errors"

//goland:noinspection GoUnusedGlobalVariable
var (
	ErrUserExists        = errors.New("username or email is already taken")
	ErrLoginFailed       = errors.New("there is an error in the login name or password")
	ErrEmailAuthRequired = errors.New("target user required email authentication")
	Err2faRequired       = errors.New("target user required two factor authentication")
	ErrNoSuchUser        = errors.New("target user was not found")
	ErrUserGone          = errors.New("target user was gone")
	ErrUserSuspended     = errors.New("target user is suspended")
	ErrInvalidPassword   = errors.New("invalid password")
	Err2faAlreadyEnabled = errors.New("target user 2fa is already enabled")
	ErrInvalid2faToken   = errors.New("invalid 2fa token")
	ErrCantEnable2fa     = errors.New("cannot enable target user's 2fa")
	ErrCantDisable2fa    = errors.New("cannot disable target user's 2fa")
	ErrInvalidKeyId      = errors.New("invalid keyId")
	ErrInvalidKey        = errors.New("invalid public key")
)
