package testredisutils

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func NewMockServer(t *testing.T, password string) *miniredis.Miniredis {
	s := miniredis.RunT(t)
	s.RequireAuth(password)
	return s
}
