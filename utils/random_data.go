package utils

import (
	"math/rand"
	"time"
)

var alphabet = []rune("abcdefghijklmnopqrstuvwxyz1234567890_-")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func RandomUsername() string {
	return RandomString(8)
}

func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}

func RandomPassword() string {
	return RandomString(12)
}
