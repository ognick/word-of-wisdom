package pow

import (
	"crypto/rand"

	"github.com/google/wire"
)

type Complexity byte

type Generator struct {
	*ProofOfWork
}

func NewGenerator(complexity Complexity) *Generator {
	return &Generator{
		ProofOfWork: NewProofOfWork(byte(complexity)),
	}
}

func (g *Generator) Generate() ([]byte, error) {
	challenge := make([]byte, ChallengeLen)
	challenge[0] = g.complexity
	if _, err := rand.Read(challenge[1:]); err != nil {
		return nil, err
	}

	return challenge, nil
}

var Set = wire.NewSet(NewGenerator)
