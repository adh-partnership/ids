package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/faa/nasr"
	"github.com/adh-partnership/ids/backend/pkg/utils"
)

func main() {
	fmt.Printf("Is today a cycle date? %+v\n", nasr.IsDateNewCycle(utils.Now()))
	fmt.Printf("Next cycle date? %+v\n", nasr.GetNextCycle(utils.Now()))
	fmt.Printf("Current cycle date? %+v\n", nasr.GetCurrentCycle(utils.Now()))

	PrintMemUsage()

	apts, err := nasr.ProcessAirports()
	if err != nil {
		logger.Logger.WithField("component", "main").Errorf("unable to process airports: %s", err)
		return
	}

	PrintMemUsage()
	logger.Logger.WithField("component", "main").Infof("Sleeping")
	logger.Logger.WithField("component", "main").Infof("%+v", len(apts))

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
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
