package service

type WisdomRepository interface {
	GetWisdom() string
}

type WisdomService struct {
	wisdomRepository WisdomRepository
}

func NewWisdomService(wisdomRepository WisdomRepository) *WisdomService {
	return &WisdomService{wisdomRepository: wisdomRepository}
}

func (w *WisdomService) GetWisdom() string {
	return w.wisdomRepository.GetWisdom()
}
