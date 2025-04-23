package main

import (
	"github.com/cristalhq/acmd"
)

var runner *acmd.Runner

func createRunner() (runner *acmd.Runner) {
	runner = acmd.RunnerOf([]acmd.Command{
		{
			Name:        "fetch",
			Alias:       "f",
			Description: "Fetches a players beatleader replays and corresponding leaderboards",
			ExecFunc:    handleFetchCmd,
		},
		{
			Name:        "generate",
			Alias:       "g",
			Description: "Generates something based on the provided arguments",
			Subcommands: []acmd.Command{
				{
					Name:        "jd-config",
					Alias:       "jd",
					Description: "Generates a jd config based on the provided players scores",
					ExecFunc:    handleJDGenCmd,
				},
			},
		},
	}, acmd.Config{
		AppName:         "beatsaber-replay-analyzer",
		AppDescription:  "Built to analyze player replays and extract information",
		Version:         "v1",
		AllowNoCommands: false,
	})

	return
}
