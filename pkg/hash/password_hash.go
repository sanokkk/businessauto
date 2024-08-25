package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

type hasher interface {
	hash(password string) string
}

type Sha256Hash struct{}

func (hasher Sha256Hash) hash(password string) string {
	sh := sha256.Sum256([]byte(password))
	hashStr := hex.EncodeToString(sh[:])

	return hashStr
}

func HashPassword(hasher hasher, password string) string {
	return hasher.hash(password)
}

type BcryptHash struct{}

func (b BcryptHash) hash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes)
}
