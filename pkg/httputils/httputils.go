package httputils

import (
	"encoding/json"
	"fmt"
	"io"
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

type FetchResponse struct {
	err      error
	Response *http.Response
}

func (r *FetchResponse) Body() (res []byte, err error) {
	if r.err != nil {
		err = r.err
		return
	}
	res, err = io.ReadAll(r.Response.Body)
	return
}

func (r *FetchResponse) BodyString() (res string, err error) {
	b, err := r.Body()
	if err == nil {
		res = string(b)
	}
	return
}

func (r *FetchResponse) BodyJSON(res interface{}) (err error) {
	b, err := r.Body()
	if err == nil {
		err = json.Unmarshal(b, res)
	}
	return
}

func Fetch(url string) (res *FetchResponse) {
	res = &FetchResponse{}
	//#nosec G107 -- implementation of generic utility
	response, err := http.Get(url)
	res.Response = response
	res.err = err
	return
}
