package airports

import (
	"errors"
	"time"

	"github.com/adh-partnership/ids/backend/internal/cache"
	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/adh-partnership/ids/backend/pkg/utils"
	"github.com/adh-partnership/ids/backend/pkg/weather"
	"gorm.io/gorm"
)

const (
	cachePrefix = "/airports"
	cacheAll    = cachePrefix + "/_all"
)

var (
	ErrInvalidAirport = errors.New("invalid airport")
)

type AirportService struct {
	db    *gorm.DB
	cache *cache.Cache
	exp   time.Duration
	hooks []func(old Airport, new Airport)
}

func NewAirportService(db *gorm.DB, c *cache.Cache) *AirportService {
	exp := time.Duration(config.GetConfig().Cache.DefaultExpiration.Airports) * time.Second

	return &AirportService{db: db, cache: c, exp: exp}
}

func (s *AirportService) AddHook(h func(old Airport, new Airport)) {
	s.hooks = append(s.hooks, h)
}

func (s *AirportService) callHooks(old Airport, new Airport) {
	for _, h := range s.hooks {
		go h(old, new)
	}
}

func (s *AirportService) GetAirports() ([]Airport, error) {
	var airports []Airport
	if airport, err := s.cache.Get(cacheAll); err == nil {
		airports = airport.([]Airport)
	} else if err == cache.ErrorKeyNotFound {
		if err := s.db.Find(&airports).Error; err != nil {
			return nil, err
		}
		if err := s.cache.Set(cacheAll, airports, s.exp); err != nil {
			return nil, err
		}
	}

	return airports, nil
}

func (s *AirportService) GetAirport(id string) (Airport, error) {
	var airport Airport
	apt, err := s.cache.Get(cachePrefix + "/" + id)
	if err == nil {
		airport = apt.(Airport)
	} else if err == cache.ErrorKeyNotFound {
		if err := s.db.Model(Airport{}).Where(Airport{FAAID: id}).Or(Airport{ICAOID: id}).First(&airport).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return Airport{}, ErrInvalidAirport
			}
		}
		if err := s.cache.Set(cachePrefix+"/"+airport.FAAID, airport, s.exp); err != nil {
			return Airport{}, err
		}
		if err := s.cache.Set(cachePrefix+"/"+airport.ICAOID, airport, s.exp); err != nil {
			return Airport{}, err
		}
	}

	return airport, nil
}

func (s *AirportService) CreateAirport(airport *Airport) error {
	if err := s.db.Create(airport).Error; err != nil {
		return err
	}
	if err := s.cache.Set(cachePrefix+"/"+airport.FAAID, airport, s.exp); err != nil {
		return err
	}
	if err := s.cache.Set(cachePrefix+"/"+airport.ICAOID, airport, s.exp); err != nil {
		return err
	}
	if err := s.cache.Delete(cacheAll); err != nil {
		return err
	}

	s.callHooks(Airport{}, *airport)

	return nil
}

func (s *AirportService) UpdateAirport(airport Airport) error {
	var oldAirport Airport
	oldAirport, err := s.GetAirport(airport.FAAID)
	if err != nil {
		return err
	}

	if oldAirport.ATIS != airport.ATIS && airport.ATIS != "" {
		airport.ATISTime = utils.Now()
	}

	if oldAirport.ArrivalATIS != airport.ArrivalATIS && airport.ArrivalATIS != "" {
		airport.ArrivalATISTime = utils.Now()
	}

	if err := s.db.Save(airport).Error; err != nil {
		return err
	}
	if err := s.cache.Set(cachePrefix+"/"+airport.FAAID, airport, s.exp); err != nil {
		return err
	}
	if err := s.cache.Set(cachePrefix+"/"+airport.ICAOID, airport, s.exp); err != nil {
		return err
	}
	if err := s.cache.Delete(cacheAll); err != nil {
		return err
	}
	s.callHooks(oldAirport, airport)

	return nil
}

func (s *AirportService) DeleteAirport(id string) error {
	var airport Airport
	if err := s.db.Where(Airport{FAAID: id}).Or(Airport{ICAOID: id}).First(&airport).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&airport).Error; err != nil {
		return err
	}
	if err := s.cache.Delete(cachePrefix + "/" + airport.FAAID); err != nil {
		return err
	}
	if err := s.cache.Delete(cachePrefix + "/" + airport.ICAOID); err != nil {
		return err
	}
	if err := s.cache.Delete(cacheAll); err != nil {
		return err
	}

	s.callHooks(airport, Airport{})

	return nil
}

func (s *AirportService) CleanupATIS() error {
	airports, err := s.GetAirports()
	if err != nil {
		return err
	}

	var errs []error

	for _, airport := range airports {
		// Check if airport.ATISTime is more than 2 hours ago
		if airport.ATISTime != nil {
			if airport.ATISTime.Before(time.Now().Add(-2 * time.Hour)) {
				airport.ATIS = ""
				airport.ATISTime = &time.Time{}
				airport.ArrivalATIS = ""
				airport.ArrivalATISTime = &time.Time{}
				airport.ArrivalRunways = ""
				airport.DepartureRunways = ""
				if err := s.UpdateAirport(airport); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	if len(errs) > 0 {
		var errStr string
		for _, err := range errs {
			errStr += err.Error() + " // "
		}
		return errors.New(errStr)
	}

	return nil
}

func (s *AirportService) UpdateWeather() error {
	err := weather.UpdateWeatherCache()
	if err != nil {
		logger.ZL.Err(err).Msg("error updating weather cache")
		return err
	}

	airports, err := s.GetAirports()
	if err != nil {
		logger.ZL.Err(err).Msg("error getting airports")
		return err
	}

	for _, airport := range airports {
		changed := false
		wx, err := weather.GetWeather(airport.ICAOID)
		if err != nil && !errors.Is(err, weather.ErrorNoWeather) {
			logger.ZL.Err(err).Msgf("error getting weather for %s, skipping", airport.FAAID)
			continue
		}

		// If we don't have weather, treat as empty so we can blank out our cache
		if errors.Is(err, weather.ErrorNoWeather) {
			wx = &weather.Weather{
				METAR: "",
				TAF:   "",
			}
		}

		if airport.METAR != wx.METAR {
			logger.ZL.Debug().Msgf("updating METAR for %s to %s", airport.FAAID, wx.METAR)
			airport.METAR = wx.METAR
			changed = true
		}
		if airport.TAF != wx.TAF {
			logger.ZL.Debug().Msgf("updating TAF for %s to %s", airport.FAAID, wx.TAF)
			airport.TAF = wx.TAF
			changed = true
		}

		if changed {
			logger.ZL.Debug().Msgf("updating airport %s", airport.FAAID)
			if err := s.UpdateAirport(airport); err != nil {
				logger.ZL.Err(err).Msgf("error updating airport %s", airport.FAAID)
				return err
			}
		}
	}

	return nil
}
