package update

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adh-partnership/ids/backend/internal/cache"
	"github.com/adh-partnership/ids/backend/internal/database"
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/charts"
	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/faa/nasr"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var log zerolog.Logger

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
			log = logger.ZL.With().Str("component", "update").Logger()
			log.Info().Msg("Updating Data...")
			log.Info().Msgf("config=%s", c.String("config"))
			log.Info().Msgf("airports=%s", c.String("airports"))

			log.Info().Msg("Parsing config...")
			config.ParseConfig(c.String("config"))

			log.Info().Msg("Parsing airports...")
			config.ParseAirports(c.String("airports"))

			log.Info().Msg("Configuring Database Repository")
			db, err := database.New(database.DBOptions{
				Host:     config.GetConfig().Database.Host,
				Port:     fmt.Sprint(config.GetConfig().Database.Port),
				User:     config.GetConfig().Database.Username,
				Password: config.GetConfig().Database.Password,
				Database: config.GetConfig().Database.DatabaseName,
				Driver:   config.GetConfig().Database.Driver,
			})
			if err != nil {
				log.Error().Msgf("unable to configure database: %s", err)
				return err
			}

			if config.GetConfig().Database.AutoMigrate {
				log.Info().Msg("Running migrations")
				err := db.Migrate(
					airports.Airport{},
					charts.Chart{},
				)
				if err != nil {
					log.Error().Msgf("unable to run migrations: %s", err)
					return err
				}
			}

			log.Info().Msg("Creating Cache Service")
			// This doesn't matter much... we won't be around long enough
			che := cache.NewCache(context.Background(), nil, 5*time.Minute, 10*time.Minute)

			log.Info().Msg("Creating Airport Service")
			aptService := airports.NewAirportService(db.DB, che)

			apts, err := nasr.ProcessAirports()
			if err != nil {
				log.Error().Msgf("unable to process airports: %s", err)
				return err
			}
			log.Info().Msgf("Found %d airports", len(apts))

			for _, apt_id := range config.GetAirports() {
				log.Info().Msgf("Processing airport %s", apt_id)
				apt, ok := apts[apt_id]
				if !ok {
					log.Error().Msgf("unable to find airport %s", apt_id)
					continue
				}

				log.Info().Msgf("Airport %s: %+v", apt_id, apt)
				arpt, err := aptService.GetAirport(apt_id)

				magVar := apt.MagnaticVariation
				if apt.MagneticHemisphere == "W" {
					magVar = -1 * magVar
				}

				if errors.Is(err, airports.ErrInvalidAirport) {
					aptModel := airports.Airport{
						FAAID:  apt_id,
						ICAOID: apt.ICAO,
						MagVar: magVar,
					}
					err = aptService.CreateAirport(&aptModel)
					if err != nil {
						log.Error().Msgf("unable to create airport %s: %s", apt_id, err)
						continue
					}
				} else {
					arpt.MagVar = magVar
					err = aptService.UpdateAirport(arpt)
					if err != nil {
						log.Error().Msgf("unable to update airport %s: %s", apt_id, err)
						continue
					}
				}
			}

			log.Info().Msg("Done")

			return nil
		},
	}
}
