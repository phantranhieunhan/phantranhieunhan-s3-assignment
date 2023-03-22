package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SuccessRes struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Paging  interface{} `json:"paging,omitempty"`
	Filter  interface{} `json:"filter,omitempty"`
	Custom  interface{} `json:"-"`
}

func (p SuccessRes) MarshalJSON() ([]byte, error) {
	// Turn p into a map
	type SuccessRes_ SuccessRes // prevent recursion
	b, _ := json.Marshal(SuccessRes_(p))

	var m map[string]json.RawMessage
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	if p.Custom != nil {
		c, err := json.Marshal(p.Custom)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(c, &m)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(m)
}

func NewSuccessResponse(data, paging, filter interface{}, custom interface{}) *SuccessRes {
	return &SuccessRes{
		Success: true,
		Data:    data,
		Paging:  paging,
		Filter:  filter,
		Custom:  custom,
	}
}

func SimpleSuccessResponse(data interface{}) *SuccessRes {
	return NewSuccessResponse(data, nil, nil, nil)
}

func CustomSuccessResponse(custom interface{}) *SuccessRes {
	return NewSuccessResponse(nil, nil, nil, custom)
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
