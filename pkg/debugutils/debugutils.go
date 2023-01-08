package debugutils

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func PrintJSON(val interface{}) {
	res, err := json.MarshalIndent(val, "> ", "  ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}
}

func FileLine() (s string) {
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		s = fmt.Sprintf("%s:%d", fileName, fileLine)
	}
	return
}
