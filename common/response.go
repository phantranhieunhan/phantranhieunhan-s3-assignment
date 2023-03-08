package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SuccessRes struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data,omitempty"`
	Paging  interface{}            `json:"paging,omitempty"`
	Filter  interface{}            `json:"filter,omitempty"`
	Custom  map[string]interface{} `json:"-"`
}

func (p SuccessRes) MarshalJSON() ([]byte, error) {
	// Turn p into a map
	type SuccessRes_ SuccessRes // prevent recursion
	b, _ := json.Marshal(SuccessRes_(p))

	var m map[string]json.RawMessage
	_ = json.Unmarshal(b, &m)

	// Add tags to the map, possibly overriding struct fields
	for k, v := range p.Custom {
		// if overriding struct fields is not acceptable:
		// if _, ok := m[k]; ok { continue }
		b, _ = json.Marshal(v)
		m[k] = b
	}

	return json.Marshal(m)
}

func NewSuccessResponse(data, paging, filter interface{}, custom map[string]interface{}) *SuccessRes {
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

func CustomSuccessResponse(custom map[string]interface{}) *SuccessRes {
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
