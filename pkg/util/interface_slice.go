package util

import (
	"errors"
	"reflect"
)

func InterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	emptySlice := make([]interface{}, 0)
	if s.Kind() != reflect.Slice {
		return emptySlice, errors.New("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return emptySlice, errors.New("Slice is nil")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}
