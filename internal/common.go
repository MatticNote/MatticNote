package internal

import (
	"crypto/rand"
	"math"
	"math/big"
	mr "math/rand"
)

const (
	MinInternalIDString = "000000000000000000000000000"
	MaxInternalIDString = "aWgEPTl1tmebfsQzFP4bxwgy80V"
)

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRandString(n int) string {
	seed, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}

	mr.Seed(seed.Int64())

	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[mr.Intn(len(characterRunes))]
	}

	return string(b)
}
