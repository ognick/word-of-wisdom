package main

import (
	"word_of_wisdom/internal/server/app"
)

func main() {
	application, err := app.InitApp()
	if err != nil {
		panic(err)
	}
	application.Run()
}
