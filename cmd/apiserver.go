package main

import (
	"math/rand"
	"os"
	"s8s/cmd/kube-apiserver/app"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())


	command := app.NewAPIServerCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
