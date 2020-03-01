package utils

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)


func Hash(host string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(host))
	return hex.EncodeToString(h.Sum(nil))
}
