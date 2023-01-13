package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
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

func ContainsObjectID(s []primitive.ObjectID, str primitive.ObjectID) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GeneratePin(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// NewSHA256 ...
func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
