package starnametesttools

import (
	"math/rand"
	"strings"
	"time"
)

func randomString(length int) string {
	rand.NewSource(time.Now().UnixNano())

	var alphabet string = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder

	l := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return sb.String()
}
