package pow

const (
	ChallengeLen = 8
	SolutionLen  = 8
)

type ProofOfWork struct {
	complexity byte
	mask       []byte
}

func NewProofOfWork(complexity byte) *ProofOfWork {
	mask := make([]byte, 0)
	for maskLen := int(complexity); maskLen > 0; maskLen -= 8 {
		shift := 8
		if maskLen < shift {
			shift = maskLen
		}
		mask = append(mask, byte((0xFF<<(8-shift))&0xFF))
	}

	return &ProofOfWork{
		complexity: complexity,
		mask:       mask,
	}
}

func (pow *ProofOfWork) isMaskMatched(hash []byte) bool {
	if len(pow.mask) > len(hash) {
		return false
	}

	for i, m := range pow.mask {
		h := hash[i]
		if m&h != 0 {
			return false
		}
	}

	return true
}

func (pow *ProofOfWork) Validate(challenge, solution []byte) bool {
	if len(solution) != SolutionLen {
		return false
	}

	hash := Hash(challenge, solution)
	return pow.isMaskMatched(hash)
}
