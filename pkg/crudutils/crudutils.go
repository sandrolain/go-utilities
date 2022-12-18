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
