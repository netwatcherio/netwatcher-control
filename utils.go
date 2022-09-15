package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"strings"
)

func toRawJson(v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&v)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// NewSHA256 ...
func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
