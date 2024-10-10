package pow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	maxComplexity = 10
)

func Test_Positive(t *testing.T) {
	generator := NewGenerator(5)
	challenge, err := generator.Generate()
	require.NoError(t, err)

	solver := NewSolver()
	solution, err := solver.Solve(challenge)
	require.NoError(t, err)

	ok := generator.Validate(challenge, solution)
	require.True(t, ok)
}

func Test_BadSolution(t *testing.T) {
	generator := NewGenerator(5)
	challenge, err := generator.Generate()
	require.NoError(t, err)

	ok := generator.Validate(challenge, make([]byte, SolutionLen))
	require.False(t, ok)

	ok = generator.Validate(challenge, make([]byte, 0))
	require.False(t, ok)
}

func Test_Generate_Solve_Validate(t *testing.T) {
	generator := NewGenerator(5)
	for complexity := byte(1); complexity <= maxComplexity; complexity++ {
		t.Run(fmt.Sprintf("complexity_%d", complexity), func(t *testing.T) {
			challenge, err := generator.Generate()
			require.NoError(t, err)

			solver := NewSolver()
			solution, err := solver.Solve(challenge)
			require.NoError(t, err)
			ok := generator.Validate(challenge, solution)
			require.True(t, ok)
		})
	}
}

func Benchmark_Generate_Solve(b *testing.B) {
	generator := NewGenerator(5)
	for complexity := byte(1); complexity <= maxComplexity; complexity++ {
		b.Run(fmt.Sprintf("complexity_%d", complexity), func(b *testing.B) {
			challenge, err := generator.Generate()
			require.NoError(b, err)

			solver := NewSolver()
			_, err = solver.Solve(challenge)
			require.NoError(b, err)
		})
	}
}
