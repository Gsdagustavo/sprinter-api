package derr

var (
	BadRequestError     = NewClientError("BAD_REQUEST", "bad request")
	UnauthorizedError   = NewClientError("UNAUTHORIZED", "unauthorized")
	NotFoundError       = NewClientError("NOT FOUND", "not found")
	InternalServerError = NewRepositoryError("INTERNAL_SERVER_ERROR", "internal server error")
)

var (
	InvalidParameterError = NewClientError("INVALID_PARAMETER", "invalid parameter")
)

var (
	InvalidEmail       = NewClientError("INVALID_EMAIL", "invalid email")
	InvalidPassword    = NewClientError("INVALID_PASSWORD", "invalid password")
	InvalidCredentials = NewClientError("INVALID_CREDENTIALS", "invalid credentials")
)
