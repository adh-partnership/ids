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
			&cli.StringFlag{
				Name:    "log-format",
				Value:   "text",
				Usage:   "Set the logging format",
				EnvVars: []string{"LOG_FORMAT"},
				Aliases: []string{"f"},
			},
		},
		Before: func(c *cli.Context) error {
			format := c.String("log-format")
			logger.New(format)

			return nil
		},
	}
}
