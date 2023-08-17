package utils

import (
	"golang.org/x/crypto/bcrypt"
)

var (
	PASS_LEN int = 6
)

func CompareHashAndPassword(f string, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(f), []byte(s))
}

func GetHash(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p), 14)
}
