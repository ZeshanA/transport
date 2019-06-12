package main

import (
	"detector/api"
	"detector/eval"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please pass evaluation mode (-e) or server mode (-s) as a CLI argument")
	}
	mode := os.Args[1]
	if mode == "-e" {
		eval.Evaluate()
	} else {
		api.Start()
	}
}
