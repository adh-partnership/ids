package server

import (
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/urfave/cli/v2"
)

var log = logger.ZL.With().Str("component", "server").Logger()

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start backend server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "config.yaml",
				Usage:   "Load configuration from `FILE`",
				Aliases: []string{"c"},
				EnvVars: []string{"CONFIG"},
			},
		},
		Action: func(c *cli.Context) error {
			log.Info().Msgf("Starting server...")
			log.Info().Msgf("config=%s", c.String("config"))

			return nil
		},
	}
}
