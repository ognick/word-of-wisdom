package service

import "word_of_wisdom/pkg/pow"

type ChallengeService struct {
	complexity byte
	mask       []byte
}

func NewChallengeService(complexity byte) *ChallengeService {
	return &ChallengeService{
		mask: pow.MakeValidationMask(complexity),
	}
}

func (c *ChallengeService) GenerateChallenge() ([]byte, error) {
	return pow.Generate(c.complexity)
}

func (c *ChallengeService) ValidateSolution(challenge, solution []byte) bool {
	return pow.Validate(challenge, solution, c.mask)
}
