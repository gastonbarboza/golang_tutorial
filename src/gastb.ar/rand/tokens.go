package rand

// The rand package wraps crypto/rand and generates user remember tokens

import (
	"crypto/rand"
	"encoding/base64"
)

// Number of bytes used to generate tokens
const RememberTokenBytes = 32

// Bytes generates n random bytes using crypto/rand
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_,err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String returns a base64 URL encoded version of the byte slice of size n
// generated by Bytes
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken generates remember tokens of a predetermined byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
