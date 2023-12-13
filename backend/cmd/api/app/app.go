package app

import (
	"github.com/adh-partnership/ids/backend/cmd/api/app/server"
	"github.com/adh-partnership/ids/backend/cmd/api/app/update"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/urfave/cli/v2"
)

func NewRootCommand() *cli.App {
	return &cli.App{
		Name:  "app",
		Usage: "PAZA Information Display Service Backend",
		Commands: []*cli.Command{
			server.NewCommand(),
			update.NewCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "Set the logging level",
				EnvVars: []string{"LOG_LEVEL"},
				Aliases: []string{"l"},
			},
		},
		Before: func(c *cli.Context) error {
			logger.New(c.String("log-level"))

			return nil
		},
	}
}
