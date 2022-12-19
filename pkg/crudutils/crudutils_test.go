package crudutils

import (
	"testing"
)

func TestErrorTypes(t *testing.T) {
	{
		err := NotFound("")
		if !IsNotFound(err) {
			t.Fatalf("Error is not NotFound: %v", err)
		}
	}

	{
		err := NotAuthorized("")
		if !IsNotAuthorized(err) {
			t.Fatalf("Error is not NotAuthorized: %v", err)
		}
	}

	{
		err := NotFound("")
		if IsNotAuthorized(err) {
			t.Fatalf("Error detected as wrong type: %v", err)
		}
	}
}
