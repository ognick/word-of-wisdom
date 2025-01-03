package main

import (
	"word_of_wisdom/internal/server/app"
)

func main() {
	application, err := app.InitializeApp()
	if err != nil {
		panic(err)
	}
	application.Run()
}
