package errorutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithError(t *testing.T) {
	count := 0
	counter := WithError[string](func(err error) {
		count++
	})

	noErrorFn := func() (string, error) {
		return "hello", nil
	}

	errorFn := func() (string, error) {
		return "", fmt.Errorf("this is an error")
	}

	res := counter(noErrorFn())
	assert.Equal(t, count, 0)
	assert.Equal(t, res, "hello")

	res = counter(errorFn())
	assert.Equal(t, count, 1)
	assert.Equal(t, res, "")
}
