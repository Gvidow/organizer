package service

import (
	"crypto/rand"
	"fmt"
)

const SizeIdSession = 64

func NewIdSession() string {
	id := make([]byte, SizeIdSession/2)
	rand.Read(id)
	return fmt.Sprintf("%x", id)
}
