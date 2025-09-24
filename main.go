package main

func main() {

	runner = createRunner()

	if err := runner.Run(); err != nil {
		runner.Exit(err)
	}

}
