package main

import (
	"runtime"
	"time"

	"github.com/adh-partnership/ids/backend/pkg/faa/nasr"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/utils"
)

func main() {
	logger.New("info")
	logger.ZL.Info().Msgf("Is today a cycle date? %+v", nasr.IsDateNewCycle(utils.Now()))
	logger.ZL.Info().Msgf("Next cycle date? %+v", nasr.GetNextCycle(utils.Now()))
	logger.ZL.Info().Msgf("Current cycle date? %+v", nasr.GetCurrentCycle(utils.Now()))

	PrintMemUsage()

	apts, err := nasr.ProcessAirports()
	if err != nil {
		logger.ZL.Error().Msgf("unable to process airports: %s", err)
		return
	}

	PrintMemUsage()
	logger.ZL.Info().Msg("Sleeping")
	logger.ZL.Info().Msgf("%+v", len(apts))

	runtime.GC()
	PrintMemUsage()

	for {
		// sleep for a minute
		time.Sleep(1 * time.Minute)
		PrintMemUsage()
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logger.ZL.Info().Msgf("Alloc = %v MiB", m.Alloc/1024/1024)
	logger.ZL.Info().Msgf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	logger.ZL.Info().Msgf("\tSys = %v MiB", m.Sys/1024/1024)
	logger.ZL.Info().Msgf("\tNumGC = %v", m.NumGC)
}
