package envutils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func GetEnvString(name string, def string) string {
	valStr := os.Getenv(name)
	if valStr != "" {
		return valStr
	}
	return def
}

func RequireEnvString(name string) (string, error) {
	valStr := GetEnvString(name, "")
	if valStr == "" {
		return "", fmt.Errorf("empty required environment variable \"%v\"", name)
	}
	return valStr, nil
}

func RequireEnvBase64(name string) ([]byte, error) {
	valStr, err := RequireEnvString(name)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(valStr)
}

func RequireEnvPath(name string) (res string, err error) {
	if res, err = RequireEnvString(name); err != nil {
		return
	}
	if !filepath.IsAbs(res) {
		err = fmt.Errorf("not a valid path: %v", res)
		return
	}
	_, err = os.Stat(res)
	return
}

func GetEnvInt(name string, def int) (int, error) {
	valStr := os.Getenv(name)
	if valStr != "" {
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return 0, err
		}
		return val, nil
	}
	return def, nil
}
