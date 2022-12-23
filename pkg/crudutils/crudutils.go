package crudutils

import "fmt"

func formatMessageByValue(msg string, value string) string {
	if value == "" {
		return msg
	}
	return fmt.Sprintf("%v: %v", msg, value)
}

type NotFoundError struct {
	value string
}

func (m *NotFoundError) Error() string {
	return formatMessageByValue("Not Found", m.value)
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
	return formatMessageByValue("Not Authorized", m.value)
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
	return formatMessageByValue("Invalid Value", m.value)
}

func InvalidValue(value string) error {
	return &InvalidValueError{value}
}

func IsInvalidValue(e error) bool {
	_, ok := e.(*InvalidValueError)
	return ok
}

type ExpiredResourceError struct {
	value string
}

func (m *ExpiredResourceError) Error() string {
	return formatMessageByValue("Expired Resource", m.value)
}

func ExpiredResource(value string) error {
	return &ExpiredResourceError{value}
}

func IsExpiredResource(e error) bool {
	_, ok := e.(*ExpiredResourceError)
	return ok
}
