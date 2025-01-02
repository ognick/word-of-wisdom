package wisdom

type wisdomRepository interface {
	GetWisdom() string
}

type DepRepos struct {
	Wisdom wisdomRepository
}

type Usecase struct {
	repos DepRepos
}

func NewService(repos DepRepos) *Usecase {
	return &Usecase{repos: repos}
}

func (u *Usecase) GetWisdom() string {
	return u.repos.Wisdom.GetWisdom()
}
