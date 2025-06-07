package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.errors) == 0 {
		return ""
	}
	var b strings.Builder

	b.WriteString(strconv.Itoa(len(e.errors)))
	b.WriteString(" errors occured:\n")
	for _, err := range e.errors {
		b.WriteString("\t* ")
		b.WriteString(err.Error())
	}
	b.WriteString("\n")
	return b.String()
}

func Append(err error, errs ...error) *MultiError {
	var result *MultiError

	if err == nil {
		result = &MultiError{}
	} else {
		if customErr, ok := err.(*MultiError); ok && customErr != nil {
			result = &MultiError{
				errors: append([]error{}, customErr.errors...),
			}
		} else {
			result = &MultiError{errors: []error{err}}
		}
	}

	result.errors = append(result.errors, errs...)

	return result
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
