package cryptoutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func RandomBytes(length int) (res []byte, err error) {
	res = make([]byte, length)
	_, err = rand.Read(res)
	return
}

func RandomBytesBase64(length int) (res string, err error) {
	b, err := RandomBytes(length)
	if err != nil {
		return
	}
	res = base64.StdEncoding.EncodeToString(b)
	return
}

func BcryptHash(value []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(value, 14)
}

func BcryptCompare(value []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, value)
	return err == nil
}

func Sha256Hash(value []byte) []byte {
	h := sha256.New()
	h.Write(value)
	return h.Sum(nil)
}

func Sha256Compare(value []byte, hash []byte) bool {
	h := Sha256Hash(value)
	return bytes.Equal(hash, h)
}

func Encrypt(value []byte, passPhrase [32]byte) ([]byte, error) {
	block, err := aes.NewCipher(passPhrase[0:])
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	return aesGCM.Seal(nonce, nonce, value, nil), nil
}

func EncryptWithHash(value []byte, passPhrase [32]byte) ([]byte, []byte, error) {
	hash := Sha256Hash(value)
	enc, err := Encrypt(value, passPhrase)
	if err != nil {
		return nil, nil, err
	}
	return enc, hash, nil
}

func Decrypt(value []byte, passPhrase [32]byte) ([]byte, error) {
	block, err := aes.NewCipher(passPhrase[0:])
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, secValue := value[:nonceSize], value[nonceSize:]
	return aesGCM.Open(nil, nonce, secValue, nil)
}

func DecryptAndVerify(value []byte, passPhrase [32]byte, hash []byte) ([]byte, error) {
	dec, err := Decrypt(value, passPhrase)
	if err != nil {
		return nil, err
	}
	if !Sha256Compare(dec, hash) {
		return nil, fmt.Errorf("Decrypted value not match hash")
	}
	return dec, nil
}
