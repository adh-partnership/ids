package jobs

import (
	"errors"

	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"github.com/adh-partnership/ids/backend/internal/domain/pireps"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/go-co-op/gocron/v2"
)

const (
	JOB_PIREP_EXPIRY   = "pirep_expiry"
	JOB_ATIS_CLEANUP   = "atis_cleanup"
	JOB_UPDATE_WEATHER = "update_weather"
)

type JobManager struct {
	airportService *airports.AirportService
	pirepsService  *pireps.PIREPService
	jobs           map[string]*gocron.Job
	scheduler      gocron.Scheduler
}

func New(airportService *airports.AirportService, pirepsService *pireps.PIREPService) (*JobManager, error) {
	m := &JobManager{
		airportService: airportService,
		pirepsService:  pirepsService,
		jobs:           make(map[string]*gocron.Job),
	}

	sched, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	j, err := sched.NewJob(
		gocron.CronJob("*/5 * * * *", false),
		gocron.NewTask(
			func() {
				logger.ZL.Info().Msg("Deleting expired PIREPs")
				err := m.pirepsService.DeleteExpiredPIREPs()
				if err != nil {
					logger.ZL.Error().Err(err).Msg("Error deleting expired PIREPs")
				}
			},
		),
	)
	if err != nil {
		return nil, errors.New("error creating pirep expiry job " + err.Error())
	}

	m.jobs[JOB_PIREP_EXPIRY] = &j

	j, err = sched.NewJob(
		gocron.CronJob("*/10 * * * *", false),
		gocron.NewTask(
			func() {
				logger.ZL.Info().Msg("Cleaning up ATIS")
				err := m.airportService.CleanupATIS()
				if err != nil {
					logger.ZL.Error().Err(err).Msg("Error cleaning up ATIS")
				}
			},
		),
	)
	if err != nil {
		return nil, errors.New("error creating atis cleanup job " + err.Error())
	}
	m.jobs[JOB_ATIS_CLEANUP] = &j

	j, err = sched.NewJob(
		gocron.CronJob("*/2 * * * *", false),
		gocron.NewTask(func() {
			logger.ZL.Info().Msg("Updating weather")
			err := airportService.UpdateWeather()
			if err != nil {
				logger.ZL.Error().Err(err).Msg("Error updating weather")
			}
		}),
	)
	if err != nil {
		return nil, errors.New("error creating weather update job " + err.Error())
	}
	m.jobs[JOB_UPDATE_WEATHER] = &j

	m.scheduler = sched

	return m, nil
}

func (j *JobManager) Start() {
	go j.scheduler.Start()
}

func (j *JobManager) Stop() error {
	return j.scheduler.Shutdown()
}
