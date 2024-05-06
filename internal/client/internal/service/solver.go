package service

import "word_of_wisdom/pkg/pow"

type SolverService struct {
}

func NewSolverService() *SolverService {
	return &SolverService{}
}

func (s *SolverService) Solve(challenge []byte) ([]byte, error) {
	return pow.Solve(challenge)
}
