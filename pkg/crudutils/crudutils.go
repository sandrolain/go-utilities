package crudutils

import "fmt"

type NotFoundError struct {
	value string
}

func (m *NotFoundError) Error() string {
	if m.value == "" {
		return "Not Found"
	}
	return fmt.Sprintf("Not Found: %v", m.value)
}

func NotFound(value string) error {
	return &NotFoundError{value}
}

func IsNotFound(e error) bool {
	_, ok := e.(*NotFoundError)
	return ok
}

type NotAuthorizedError struct {
	value string
}

func (m *NotAuthorizedError) Error() string {
	if m.value == "" {
		return "Not Authorized"
	}
	return fmt.Sprintf("Not Authorized: %v", m.value)
}

func NotAuthorized(value string) error {
	return &NotAuthorizedError{value}
}

func IsNotAuthorized(e error) bool {
	_, ok := e.(*NotAuthorizedError)
	return ok
}

type InvalidValueError struct {
	value string
}

func (m *InvalidValueError) Error() string {
	if m.value == "" {
		return "Invalid Value"
	}
	return fmt.Sprintf("Invalid Value: %v", m.value)
}

func InvalidValue(value string) error {
	return &InvalidValueError{value}
}

func IsInvalidValue(e error) bool {
	_, ok := e.(*InvalidValueError)
	return ok
}
