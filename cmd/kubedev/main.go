package main

import (
	"knativetest/pkg/dev/kubedev"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		args := append([]string{"help"}, os.Args[1:]...)
		kubedev.Cmd.SetArgs(args)
	}

	if err := kubedev.Cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
