package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonIterations = 1
	argonMemory     = 64 * 1024
	argonThreads    = 4
	argonKeyLength  = 32
	saltLength      = 16
)

var (
	ErrPasswordEmpty           = errors.New("password must not be empty")
	ErrInvalidPasswordEncoding = errors.New("invalid password encoding")
	ErrPasswordMismatch        = errors.New("incorrect password")
)

func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", ErrPasswordEmpty
	}

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", ErrorHandler(err, "unable to generate password salt")
	}

	hash := argon2.IDKey([]byte(plain), salt, argonIterations, argonMemory, argonThreads, argonKeyLength)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	return saltBase64 + "." + hashBase64, nil
}

func VerifyPassword(plain, encoded string) error {
	if encoded == "" {
		return ErrInvalidPasswordEncoding
	}

	parts := strings.Split(encoded, ".")
	if len(parts) != 2 {
		return ErrInvalidPasswordEncoding
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return ErrInvalidPasswordEncoding
	}

	storedHash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return ErrInvalidPasswordEncoding
	}

	derived := argon2.IDKey([]byte(plain), salt, argonIterations, argonMemory, argonThreads, argonKeyLength)
	if len(derived) != len(storedHash) {
		return ErrPasswordMismatch
	}
	if subtle.ConstantTimeCompare(derived, storedHash) != 1 {
		return ErrPasswordMismatch
	}

	return nil
}
