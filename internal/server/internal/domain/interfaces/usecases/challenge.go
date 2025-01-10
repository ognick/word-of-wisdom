package usecases

type Challenge interface {
	GenerateChallenge() ([]byte, error)
	ValidateSolution(challenge, solution []byte) bool
}
