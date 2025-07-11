package hasher

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) string
}

type SHA1Hasher struct {
	salt string
}

func NewSHA512Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (s *SHA1Hasher) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.salt)))
}
