package httputils

import (
	"fmt"
	"net/http"
	"strings"
)

func GetRequestBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 || len(splitToken[1]) == 0 {
		return "", fmt.Errorf("invalid Authorization Bearer Token \"%s\"", reqToken)
	}
	return splitToken[1], nil
}
