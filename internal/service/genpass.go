package service

import (
	"crypto/rand"
	"errors"
	"math/big"
)

//go:generate mockery --name=GenPass --filename genpass.go --outpkg svc_mocks

// GenPass generates random strings of letters and digits.
type GenPass interface {
	GeneratePassword(length int) (string, error)
}

type genPass struct{}

func NewGenPass() GenPass {
	return &genPass{}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var ErrInvalidLength = errors.New("password length must be between 6 and 64")

func (g *genPass) GeneratePassword(length int) (string, error) {
	if length < 6 || length > 64 {
		return "", ErrInvalidLength
	}

	password := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}
