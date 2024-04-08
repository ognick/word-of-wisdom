package pow

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math"
)

var NoSolutionError = errors.New("no solution")

const (
	ChallengeLen = 8
	SolutionLen  = 8
)

func GetComplexity(challenge []byte) byte {
	return challenge[0]
}

func MakeValidationMask(complexity byte) []byte {
	mask := make([]byte, 0)
	for maskLen := int(complexity); maskLen > 0; maskLen -= 8 {
		shift := 8
		if maskLen < shift {
			shift = maskLen
		}
		mask = append(mask, byte((0xFF<<(8-shift))&0xFF))
	}

	return mask
}

func IsMaskMatched(mask, hash []byte) bool {
	if len(mask) > len(hash) {
		return false
	}

	for i, m := range mask {
		h := hash[i]
		if m&h != 0 {
			return false
		}
	}

	return true
}

func Generate(complexity byte) ([]byte, error) {
	challenge := make([]byte, ChallengeLen)
	challenge[0] = complexity
	if _, err := rand.Read(challenge[1:]); err != nil {
		return nil, err
	}

	return challenge, nil
}

func Validate(challenge, solution, mask []byte) bool {
	if len(solution) != SolutionLen {
		return false
	}

	hash := Hash(challenge, solution)
	return IsMaskMatched(mask, hash)
}

func Solve(challenge []byte) ([]byte, error) {
	mask := MakeValidationMask(GetComplexity(challenge))
	for i := uint64(0); i < math.MaxUint64; i++ {
		solution := make([]byte, SolutionLen)
		binary.BigEndian.PutUint64(solution, i)
		if Validate(challenge, solution, mask) {
			return solution, nil
		}
	}

	return nil, NoSolutionError
}
