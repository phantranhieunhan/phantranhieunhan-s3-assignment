package common

import (
	"errors"
	"net/mail"
	"reflect"
)

var ErrValidateEmail = errors.New("email is not valid")

func ValidateRequired(value interface{}, fieldName string) error {
	vType := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	isErr := false
	if vType == nil {
		isErr = true
	} else {
		switch vType.Kind() {
		case reflect.String:
			if val.String() == "" {
				isErr = true
			}
		case reflect.Array, reflect.Slice:
			if val.Len() == 0 {
				isErr = true
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64:
			if val.Int() == 0 {
				isErr = true
			}

		}

	}

	if isErr {
		return ErrInvalidRequest(nil, fieldName)
	}
	return nil
}

func ValidateEmail(address string) error {
	_, err := mail.ParseAddress(address)

	if err != nil {
		return ErrInvalidRequest(ErrValidateEmail, "")
	}

	return nil

}
