package main

import (
	"os"

	"github.com/EdgarPost/envsec/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
