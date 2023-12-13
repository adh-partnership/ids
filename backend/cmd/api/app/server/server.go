package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adh-partnership/ids/backend/internal/cache"
	"github.com/adh-partnership/ids/backend/internal/database"
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/charts"
	"github.com/adh-partnership/ids/backend/internal/handlers"
	"github.com/adh-partnership/ids/backend/internal/redis"
	"github.com/adh-partnership/ids/backend/internal/server"
	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var log zerolog.Logger

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
			logger.ZL.Info().Msg("Configuring logger")
			log = logger.ZL.With().Str("component", "server").Logger()
			log.Info().Msgf("Starting server...")
			log.Info().Msgf("config=%s", c.String("config"))

			log.Info().Msg("Parsing config...")
			err := config.ParseConfig(c.String("config"))
			if err != nil {
				log.Error().Msgf("unable to parse config: %s", err)
				return err
			}

			log.Info().Msg("Connecting to database...")
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
				log.Info().Msg("Auto migrating database...")
				err = db.Migrate(
					airports.Airport{},
					charts.Chart{},
				)
				if err != nil {
					log.Error().Msgf("unable to migrate database: %s", err)
					return err
				}
			}

			var r *redis.Redis
			if config.GetConfig().Cache.Driver == "redis" {
				log.Info().Msg("Connecting to Redis...")
				r = redis.New(&config.GetConfig().Cache)
			}

			log.Info().Msg("Configuring cache service")
			che := cache.NewCache(context.TODO(), r, 15*time.Minute, 2*time.Minute)

			log.Info().Msg("Building Server service")
			s := server.New()

			log.Info().Msg("Registering services...")
			as := airports.NewAirportService(db.DB, che)
			cs := charts.NewChartService(db.DB, che, as)

			log.Info().Msg("Registering handlers...")
			s.Router.Route("/v1", func(r chi.Router) {
				handlers.RegisterHandlers(r, as, cs)
			})
			handlers.NewServiceHandlers(s.Router)

			log.Info().Msg("Listing defined routes:")
			chi.Walk(s.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
				logger.ZL.Info().Msgf(" - %s %s", method, route)
				return nil
			})

			log.Info().Msg("Starting server...")
			err = s.Start(config.GetConfig().Server.Mode, fmt.Sprintf(":%d", config.GetConfig().Server.Port))
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error().Msgf("unable to start server: %s", err)
			}

			return nil
		},
	}
}
