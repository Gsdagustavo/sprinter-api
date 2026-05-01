package derr

import (
	"errors"
)

type RepositoryError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e RepositoryError) Error() string {
	return e.Message
}

type ClientError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ClientError) Error() string {
	return e.Message
}

func NewRepositoryError(code string, message string) RepositoryError {
	return RepositoryError{
		Code:    code,
		Message: message,
	}
}

func NewClientError(code string, message string) ClientError {
	return ClientError{
		Code:    code,
		Message: message,
	}
}

func NewInternalError(message string) RepositoryError {
	return RepositoryError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: message,
	}
}

func NewBadRequestError(message string) ClientError {
	return ClientError{
		Code:    "BAD_REQUEST",
		Message: message,
	}
}

func NewNotFoundError(message string) ClientError {
	return ClientError{
		Code:    "NOT_FOUND",
		Message: message,
	}
}

func NewUnauthorizedError(message string) ClientError {
	return ClientError{
		Code:    "UNAUTHORIZED",
		Message: message,
	}
}

func JoinError(message string, err error) error {
	return errors.Join(errors.New(message), err)
}
