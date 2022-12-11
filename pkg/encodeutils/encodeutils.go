package encodeutils

import (
	"encoding/base64"
	"encoding/hex"
)

func Base64Encode(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func Base64Decode(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func HexEncode(value []byte) string {
	return hex.EncodeToString(value)
}

func HexDecode(value string) ([]byte, error) {
	return hex.DecodeString(value)
}
