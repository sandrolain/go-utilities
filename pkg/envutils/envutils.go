package envutils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
)

func GetEnvString(name string, def string) string {
	valStr := os.Getenv(name)
	if valStr != "" {
		return name
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
