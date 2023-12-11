package nasr

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	base     time.Time
	interval = 28
)

func init() {
	os.Setenv("TZ", "UTC")
	base = time.Date(
		2023, 11, 30, 0, 0, 0, 0, time.UTC,
	)
}

func IsDateNewCycle(dt *time.Time) bool {
	// Find number of days from base to dt
	days := int(dt.Sub(base).Hours() / 24)
	return days%interval == 0
}

func GetNextCycle(dt *time.Time) time.Time {
	// Find number of days from base to dt
	days := int(dt.Sub(base).Hours() / 24)
	// Find number of days to next cycle
	rem := days % interval
	next := interval - rem
	return dt.AddDate(0, 0, next)
}

func GetCurrentCycle(dt *time.Time) *time.Time {
	// Find the starting date of the current cycle
	next := GetNextCycle(dt).AddDate(0, 0, -interval)
	return &next
}

func URL(dt *time.Time, t string) string {
	return fmt.Sprintf(
		"https://nfdc.faa.gov/webContent/28DaySub/extra/%s_%s_CSV.zip",
		GetFormattedDate(dt),
		strings.ToUpper(t),
	)
}

func GetFormattedDate(dt *time.Time) string {
	// Format is DD_MMM_YYYY
	return fmt.Sprintf(
		"%02d_%s_%d",
		dt.Day(),
		dt.Month().String()[:3],
		dt.Year(),
	)
}
