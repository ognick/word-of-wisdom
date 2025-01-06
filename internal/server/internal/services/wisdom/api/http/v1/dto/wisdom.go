package dto

import "github.com/ognick/word_of_wisdom/internal/server/internal/domain/models"

type Wisdom struct {
	Message string `json:"message"`
}

func NewWisdom(w models.Wisdom) *Wisdom {
	return &Wisdom{
		Message: w.Content,
	}
}
