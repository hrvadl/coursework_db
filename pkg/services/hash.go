package services

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(password string) (string, error)
}

type Comparer interface {
	Compare(password string, hashed string) error
}

type Cryptor interface {
	Comparer
	Hasher
}

func NewCryptor() Cryptor {
	return &authenticator{}
}

type authenticator struct{}

func (a *authenticator) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *authenticator) Compare(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
