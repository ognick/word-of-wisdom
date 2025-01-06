package pow

import (
	"crypto/rand"
)

type Generator struct {
	*ProofOfWork
}

func NewGenerator(complexity byte) *Generator {
	return &Generator{
		ProofOfWork: NewProofOfWork(complexity),
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
