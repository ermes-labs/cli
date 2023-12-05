package main

import (
	"log"

	"github.com/ermes-labs/ermes-cli/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
