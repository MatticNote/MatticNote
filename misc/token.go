package misc

import "math/rand"

var tokenCharset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenToken(size uint8) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = tokenCharset[rand.Intn(len(tokenCharset))]
	}
	return string(b)
}
