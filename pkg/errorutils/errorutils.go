package errorutils

import (
	"fmt"
	"runtime"
)

func LineError(message string, args ...interface{}) error {
	message = fmt.Sprintf(message, args...)
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		return fmt.Errorf("[%s:%d] %s", fileName, fileLine, message)
	}
	return fmt.Errorf(message)
}
