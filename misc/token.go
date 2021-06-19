package misc

import "math/rand"

var tokenCharset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var backupCodeCharset = []rune("0123456789")

func GenToken(size uint8) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = tokenCharset[rand.Intn(len(tokenCharset))]
	}
	return string(b)
}

func GenBackupCode() [8]string {
	code := [8]string{}

	for i := 0; i < 8; i++ {
		b := make([]rune, 10)
		for l := range b {
			b[l] = backupCodeCharset[rand.Intn(len(backupCodeCharset))]
		}
		code[i] = string(b)
	}

	return code
}
