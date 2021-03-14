package main

import (
	"os"

	"github.com/rdeusser/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rdeusser/httpstat-monitor/cmd/httpstat-monitor/app"
	"github.com/rdeusser/httpstat-monitor/version"
)

type RootCommand struct{}

func (RootCommand) Init() cli.Command {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cmd := cli.Command{
		Name:    "httpstat-monitor",
		Desc:    "A service that monitors https://httpstat.us/503 and https://httpstat.us/200",
		Version: version.GetHumanVersion(),
	}

	cmd.AddCommands(
		&app.ServerCommand{},
	)

	return cmd
}

func (RootCommand) Run(args []string) error {
	return cli.PrintHelp
}

func main() {
	if _, err := cli.Run(&RootCommand{}); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
