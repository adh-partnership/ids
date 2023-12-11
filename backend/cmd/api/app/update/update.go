package update

import (
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/urfave/cli/v2"
)

var log = logger.Logger.WithField("component", "update")

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "config.yaml",
				Usage:   "Load configuration from `FILE`",
				Aliases: []string{"c"},
				EnvVars: []string{"CONFIG"},
			},
			&cli.StringFlag{
				Name:    "airports",
				Value:   "airports.json",
				Usage:   "Define airports to build/update from `FILE`",
				Aliases: []string{"a"},
				EnvVars: []string{"AIRPORTS"},
			},
		},
		Action: func(c *cli.Context) error {
			log.Infof("Updating Data...")
			log.Infof("config=%s", c.String("config"))
			log.Infof("airports=%s", c.String("airports"))

			return nil
		},
	}
}
