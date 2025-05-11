package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Конвертирует строковое выражение в формат SHA512
func ExpressionToSHA512(expression string) [64]byte {
	expression = strings.ReplaceAll(expression, " ", "")
	hash := sha512.Sum512([]byte(expression))
	return hash
}

// Конвертирует SHA512 в base64 URL encoding
func EncodeToString(hash [64]byte) string {
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Конвертирует base64 URL encoding в SHA512
func EncodedToSHA512(encoded string) ([64]byte, error) {
	hash, err := base64.URLEncoding.DecodeString(encoded)
	if len(hash) != 64 {
		return [64]byte{}, ErrInvalidBase64Decode
	}
	return [64]byte(hash), err
}

// Обертка для error в формат JSON
func ErrorResponse(err error) []byte {
	type Error struct {
		Error string `json:"error"`
	}
	data, _ := json.Marshal(Error{
		Error: err.Error(),
	})
	return data
}
