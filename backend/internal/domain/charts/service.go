package charts

import (
	"errors"
	"sync"

	"github.com/adh-partnership/ids/backend/internal/cache"
	"github.com/adh-partnership/ids/backend/internal/domain/airports"
	"gorm.io/gorm"
)

const (
	cachePrefix = "/charts"
)

var (
	ErrInvalidAiport = errors.New("invalid airport")
	ErrNoCharts      = errors.New("no charts found")
)

type ChartService struct {
	airportService *airports.AirportService
	db             *gorm.DB
	cache          *cache.Cache
	mut            sync.RWMutex
}

func NewChartService(db *gorm.DB, cache *cache.Cache, airportService *airports.AirportService) *ChartService {
	return &ChartService{
		db:    db,
		cache: cache,
	}
}

func (s *ChartService) GetCharts(apt_id string) ([]*Chart, error) {
	s.mut.RLock()
	defer s.mut.RUnlock()

	// We need to lookup the airport as we should be able to use both ICAO and FAA IDs,
	// but our chart-parser only inserts based on their FAA IDs.
	airport, err := s.airportService.GetAirport(apt_id)
	if err != nil {
		if errors.Is(err, airports.ErrInvalidAirport) {
			return nil, ErrInvalidAiport
		}
		return nil, err
	}

	var charts []*Chart
	if chart, err := s.cache.Get(cachePrefix + "/" + apt_id); err == nil {
		charts = chart.([]*Chart)
	} else if err != cache.ErrorKeyNotFound {
		if err := s.db.Model(Chart{}).Where(Chart{AirportID: airport.FAAID}).Find(&charts).Error; err != nil {
			return nil, err
		}
		if err := s.cache.Set(cachePrefix+"/"+apt_id, charts); err != nil {
			return nil, err
		}
	}

	return charts, nil
}
