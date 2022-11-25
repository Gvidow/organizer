package service

import (
	"crypto/sha256"
	"fmt"
)

const salt = "dsasgdfggoiugsgsdgisugig"

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
