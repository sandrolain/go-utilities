package pwdutils

import (
	"fmt"
	"net/url"

	pwdgen "github.com/sethvargo/go-password/password"
	pwdval "github.com/wagslane/go-password-validator"
)

const (
	MinEntroy = 60
)

func Validate(password string) error {
	return pwdval.Validate(password, MinEntroy)
}

func Generate(len int) (string, error) {
	dig := len / 4
	sim := len / 4
	return pwdgen.Generate(len, dig, sim, false, false)
}

func ExtractPasswordFromURI(uri string) (pwd string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	pwd, ok := u.User.Password()
	if !ok {
		err = fmt.Errorf("the URI does not contain a password")
	}
	return
}
