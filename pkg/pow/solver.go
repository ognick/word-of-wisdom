package pow

import (
	"encoding/binary"
	"errors"
	"math"
)

var NoSolutionError = errors.New("no solution")

type Solver struct {
}

func NewSolver() *Solver {
	return &Solver{}
}

func (s *Solver) Solve(challenge []byte) ([]byte, error) {
	complexity := challenge[0]
	pow := NewProofOfWork(complexity)
	for i := uint64(0); i < math.MaxUint64; i++ {
		solution := make([]byte, SolutionLen)
		binary.BigEndian.PutUint64(solution, i)
		if pow.Validate(challenge, solution) {
			return solution, nil
		}
	}

	return nil, NoSolutionError
}
