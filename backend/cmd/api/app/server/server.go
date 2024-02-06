package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adh-partnership/ids/backend/internal/cache"
	"github.com/adh-partnership/ids/backend/internal/database"
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/auth"
	"github.com/adh-partnership/ids/backend/internal/domain/charts"
	"github.com/adh-partnership/ids/backend/internal/domain/pireps"
	"github.com/adh-partnership/ids/backend/internal/handlers"
	"github.com/adh-partnership/ids/backend/internal/jobs"
	"github.com/adh-partnership/ids/backend/internal/middleware/session"
	"github.com/adh-partnership/ids/backend/internal/redis"
	"github.com/adh-partnership/ids/backend/internal/server"
	"github.com/adh-partnership/ids/backend/internal/signalr"
	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	log        zerolog.Logger
	Hub        *signalr.IDSHub
	JobManager *jobs.JobManager
)

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

			data, _ := json.MarshalIndent(config.GetConfig(), "", "  ")
			log.Debug().Msgf("config=%s", string(data))

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
					pireps.PIREP{},
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

			log.Info().Msg("Setting up sessions store")
			session.New(&config.GetConfig().Session)

			log.Info().Msg("Building Server service")
			s := server.New()

			log.Info().Msg("Registering services...")
			airportService := airports.NewAirportService(db.DB, che)
			authservice := auth.NewAuthService(session.Store, config.GetConfig().OAuth.Provider)
			chartService := charts.NewChartService(db.DB, che, airportService)
			pirepService := pireps.NewPIREPService(db.DB, che)

			log.Info().Msg("Configuring SignalR Hub...")
			s.Router.Route("/signalr", func(r chi.Router) {
				Hub, err = signalr.New(context.Background(), r, airportService, chartService, pirepService)
			})
			if err != nil {
				log.Error().Msgf("unable to configure signalr: %s", err)
				return err
			}

			log.Info().Msg("Configuring service to SignalR hooks...")
			Hub.ConfigureHooks(airportService, pirepService)

			log.Info().Msg("Registering handlers...")
			s.Router.Route("/v1", func(r chi.Router) {
				handlers.RegisterHandlers(r, &handlers.Services{
					AirportService: airportService,
					AuthService:    authservice,
					ChartService:   chartService,
					PIREPService:   pirepService,
				})
			})
			handlers.NewServiceHandler(s.Router, db.DB)

			log.Info().Msg("Listing defined routes:")
			chi.Walk(s.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
				logger.ZL.Info().Msgf(" - %s %s", method, route)
				return nil
			})

			log.Info().Msg("Configuring and starting jobs...")
			JobManager, err = jobs.New(airportService, pirepService)
			if err != nil {
				log.Error().Msgf("unable to configure jobs: %s", err)
				return err
			}
			JobManager.Start()

			log.Info().Msg("Preloading weather")
			err = airportService.UpdateWeather()
			if err != nil {
				log.Error().Msgf("unable to preload weather: %s", err)
				return err
			}

			log.Info().Msg("Preloading charts")
			_, err = chartService.GetAllCharts()
			if err != nil && !errors.Is(err, charts.ErrNoCharts) {
				log.Error().Msgf("unable to preload charts: %s", err)
				return err
			}

			log.Info().Msg("Starting server...")
			err = s.Start(config.GetConfig().Server.Mode, fmt.Sprintf("%s:%d", config.GetConfig().Server.IP, config.GetConfig().Server.Port), *JobManager)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error().Msgf("unable to start server: %s", err)
			}

			return nil
		},
	}
}
