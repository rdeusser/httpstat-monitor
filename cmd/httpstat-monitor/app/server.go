package app

import (
	"github.com/rdeusser/cli"
	"github.com/rs/zerolog/log"

	"github.com/rdeusser/httpstat-monitor/pkg/httpstat"
)

type ServerCommand struct{}

func (ServerCommand) Init() cli.Command {
	return cli.Command{
		Name: "server",
		Desc: "Starts the httpstat-monitor server",
		Flags: []cli.Flag{
			Debug,
		},
	}
}

func (ServerCommand) Run(args []string) error {
	srv := httpstat.NewServer(":8080",
		httpstat.Logger(log.Logger),
		httpstat.MonitorURL("https://httpstat.us/503"),
		httpstat.MonitorURL("https://httpstat.us/200"))

	return srv.Start()
}
