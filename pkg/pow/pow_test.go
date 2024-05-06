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
	challenge, err := Generate(5)
	require.NoError(t, err)

	solution, err := Solve(challenge)
	require.NoError(t, err)

	mask := MakeValidationMask(GetComplexity(challenge))
	ok := Validate(challenge, solution, mask)
	require.True(t, ok)
}

func Test_BadSolution(t *testing.T) {
	challenge, err := Generate(5)
	require.NoError(t, err)
	mask := MakeValidationMask(GetComplexity(challenge))

	ok := Validate(challenge, make([]byte, SolutionLen), mask)
	require.False(t, ok)

	ok = Validate(challenge, make([]byte, 0), mask)
	require.False(t, ok)
}

func Test_Generate_Solve_Validate(t *testing.T) {
	for complexity := byte(1); complexity <= maxComplexity; complexity++ {
		t.Run(fmt.Sprintf("complexity_%d", complexity), func(t *testing.T) {
			challenge, err := Generate(complexity)
			require.NoError(t, err)

			solution, err := Solve(challenge)
			require.NoError(t, err)
			mask := MakeValidationMask(GetComplexity(challenge))
			ok := Validate(challenge, solution, mask)
			require.True(t, ok)
		})
	}
}

func Benchmark_Generate_Solve(b *testing.B) {
	for complexity := byte(1); complexity <= maxComplexity; complexity++ {
		b.Run(fmt.Sprintf("complexity_%d", complexity), func(b *testing.B) {
			challenge, err := Generate(complexity)
			require.NoError(b, err)

			_, err = Solve(challenge)
			require.NoError(b, err)
		})
	}
}
