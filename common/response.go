package common

import (
	"errors"
	"net/http"
)

type SuccessRes struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Paging  interface{} `json:"paging,omitempty"`
	Filter  interface{} `json:"filter,omitempty"`
}

func NewSuccessResponse(data, paging, filter interface{}) *SuccessRes {
	return &SuccessRes{
		Success: true,
		Data:    data,
		Paging:  paging,
		Filter:  filter,
	}
}

func SimpleSuccessResponse(data interface{}) *SuccessRes {
	return NewSuccessResponse(data, nil, nil)
}

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

func NewErrorResponse(root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewUnauthorized(root error, msg, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Key:        key,
	}
}

func NewCustomError(root error, msg string, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}

	return NewErrorResponse(errors.New(msg), msg, msg, key)
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}

	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

func (e *AppError) ClearRoot() {
	e.RootErr = nil
	e.Log = ""
}
