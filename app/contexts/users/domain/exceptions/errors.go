package exceptions

import base_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"

var (
	ErrInvalidEmail   = base_exceptions.NewBadRequestError("invalid email format", "")
	ErrWeakPassword   = base_exceptions.NewBadRequestError("password must be at least 8 characters with uppercase, lowercase and digit", "")
	ErrUserNotFound   = base_exceptions.NewNotFoundError("user not found", "")
	ErrDuplicateEmail = base_exceptions.NewConflictError("email already exists", "")
	ErrInvalidUuid    = base_exceptions.NewBadRequestError("invalid user id", "")
)
