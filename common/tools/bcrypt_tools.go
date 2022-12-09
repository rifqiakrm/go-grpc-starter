// Package tools is a tool helper for th repository
package tools

import "golang.org/x/crypto/bcrypt"

const (
	cost = 10
)

// BcryptEncrypt encrypts a string.
func BcryptEncrypt(plainText string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), cost)
	return string(hashed), err
}

// BcryptVerifyHash compares hashed and plain string.
func BcryptVerifyHash(encrypted, plain string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(plain)); err != nil {
		return false
	}
	return true
}
