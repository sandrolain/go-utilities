package sliceutils

import "fmt"

func ConvertType[T interface{}](all []interface{}) ([]T, error) {
	res := make([]T, len(all))
	for i, v := range all {
		val, ok := v.(T)
		if !ok {
			return res, fmt.Errorf("item %v not convertible", i)
		}
		res[i] = val
	}
	return res, nil
}
