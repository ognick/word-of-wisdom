package wisdom

type wisdomRepository interface {
	GetWisdom() string
}

type DepRepos struct {
	Wisdom wisdomRepository
}

type Service struct {
	repos DepRepos
}

func NewService(repos DepRepos) *Service {
	return &Service{repos: repos}
}

func (w *Service) GetWisdom() string {
	return w.repos.Wisdom.GetWisdom()
}
