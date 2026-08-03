package main

import (
	"os"

	"github.com/yunkhngn/dforge/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
