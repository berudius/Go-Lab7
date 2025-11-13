package security

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

const (
	memory      uint32 = 64 * 1024
	iterations  uint32 = 3
	parallelism uint8  = 2
	keyLength   uint32 = 32
)

func CompareHash(apiKey string, storedHashHex string, saltHex string) bool {
	storedHash, err := hex.DecodeString(storedHashHex)
	if err != nil {
		return false
	}
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false
	}

	// Фаза 1: хешування sha512
	sha512Hasher := sha512.New()
	sha512Hasher.Write([]byte(apiKey))
	sha512Hash := sha512Hasher.Sum(nil)

	// Фаза 2: хешування хешу Argon2
	comparisonHash := argon2.IDKey(sha512Hash, salt, iterations, memory, parallelism, keyLength)

	return subtle.ConstantTimeCompare(storedHash, comparisonHash) == 1
}
