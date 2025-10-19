package exceptions

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e ApiError) Error() string {
	return e.Message
}

func NewBadRequestError(message, detail string) *ApiError {
	return &ApiError{
		Code:    400,
		Message: message,
		Detail:  detail,
	}
}

func NewNotFoundError(message, detail string) *ApiError {
	return &ApiError{
		Code:    404,
		Message: message,
		Detail:  detail,
	}
}

func NewConflictError(message, detail string) *ApiError {
	return &ApiError{
		Code:    409,
		Message: message,
		Detail:  detail,
	}
}

func NewInternalServerError(message, detail string) *ApiError {
	return &ApiError{
		Code:    500,
		Message: message,
		Detail:  detail,
	}
}

func NewUnauthorizedError(message, detail string) *ApiError {
	return &ApiError{
		Code:    401,
		Message: message,
		Detail:  detail,
	}
}

func NewForbiddenError(message, detail string) *ApiError {
	return &ApiError{
		Code:    403,
		Message: message,
		Detail:  detail,
	}
}

func NewValidationError(message, detail string) *ApiError {
	return &ApiError{
		Code:    422,
		Message: message,
		Detail:  detail,
	}
}

func NewServiceUnavailableError(message, detail string) *ApiError {
	return &ApiError{
		Code:    503,
		Message: message,
		Detail:  detail,
	}
}
