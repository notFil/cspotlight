package errs

import "errors"

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrNotFound              = errors.New("not found")
	ErrNoMatchOnNewPasswords = errors.New("no match on new passwords")
	ErrInvalidInput          = errors.New("invalid input")
	ErrImageTooLarge         = errors.New("image too large")
	ErrInvalidImageFormat    = errors.New("invalid image format")
	ErrMissingImage          = errors.New("missing image")
)
