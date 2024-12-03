package symmetricHash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var ErrMismatchHash = errors.New("hash doesn't match")

func GenerateSHA256(requestBody string) string {
	hash := sha256.New()
	hash.Write([]byte(requestBody))
	md := hash.Sum(nil)
	requestContentStr := hex.EncodeToString(md)
	return requestContentStr
}

func CompareSHA256(hashedString, plainString string) error {
	if hashedString == GenerateSHA256(plainString) {
		return nil
	}

	return ErrMismatchHash
}
