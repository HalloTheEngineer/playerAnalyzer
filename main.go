package main

import (
	"log"
	"replayAnalyzer/storage"
)

func main() {

	err := storage.Load()
	if err != nil {
		log.Fatal(err)
	}

	runner = createRunner()

	if err := runner.Run(); err != nil {
		runner.Exit(err)
	}

}
