package pow

import (
	"golang.org/x/crypto/argon2"
)

const (
	memory      = 1 * 1024 // 1 MB
	iterations  = 3
	keyLen      = 32
	parallelism = 1
)

func Hash(challenge, solution []byte) []byte {
	return argon2.IDKey(challenge, solution, iterations, memory, parallelism, keyLen)
}
