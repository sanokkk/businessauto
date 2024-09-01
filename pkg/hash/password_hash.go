package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

type hasher interface {
	hash(password string) string
	comparePasswordAndHash(inputPassword string, hash string) error
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
func ComparePasswordAndHash(hasher hasher, password string, hash string) bool {
	if err := hasher.comparePasswordAndHash(password, hash); err != nil {
		return false
	}

	return true
}

type BcryptHash struct{}

func (b BcryptHash) hash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes)
}

func (b BcryptHash) comparePasswordAndHash(inputPassword string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(inputPassword))
}
