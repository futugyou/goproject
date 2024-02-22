package main

import (
	"github/go-project/tour/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.execute err: %v", err)
	}

	// OrderCodeWtihFlag()
}
