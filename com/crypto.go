package com

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func HashByte(b []byte) []byte {
	h := sha256.New()

	if _, err := h.Write(b); err != nil {
		return nil
	}

	return h.Sum(nil)
}
