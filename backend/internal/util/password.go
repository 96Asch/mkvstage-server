package util

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func Validate(plaintext, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	if err != nil {
		return domain.NewInternalErr()
	}

	return nil
}
