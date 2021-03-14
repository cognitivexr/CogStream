package util

import (
	"math/rand"
)

const runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(l int) string {
	bytes := make([]byte, l)
	n := len(runes)
	for i := 0; i < l; i++ {
		bytes[i] = runes[rand.Intn(n)]
	}
	return string(bytes)
}
