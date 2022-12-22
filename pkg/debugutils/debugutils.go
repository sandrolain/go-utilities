package debugutils

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(val interface{}) {
	res, err := json.MarshalIndent(val, "> ", "  ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}
}
