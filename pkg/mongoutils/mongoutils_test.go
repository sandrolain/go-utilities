package mongoutils

import (
	"fmt"
	"testing"

	"github.com/sandrolain/go-utilities/pkg/testmongoutils"
)

func TestMain(m *testing.M) {
	testmongoutils.MockServer(m, "6.0", "user", "password")
}

func TestURI(t *testing.T) {
	fmt.Print(testmongoutils.GetMockServerURI())
}
