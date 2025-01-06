package main

import (
	"github.com/ognick/word_of_wisdom/internal/server/app"
)

func main() {
	application, err := app.InitializeApp()
	if err != nil {
		panic(err)
	}
	application.Run()
}
