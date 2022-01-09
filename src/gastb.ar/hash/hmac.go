package hash

// The hash package wraps the standard hash package and others in order
// to hash remember tokens and other strings.

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// HMAC is a wrapper around the hash.Hash interface
type HMAC struct {
	hmac hash.Hash
}

// NewHMAC creates and returns an HMAC object from a secret key 
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC {
		hmac: h,
	}
}

func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
