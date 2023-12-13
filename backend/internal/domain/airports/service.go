package airports

import (
	"errors"
	"time"

	"github.com/adh-partnership/ids/backend/internal/cache"
	"github.com/adh-partnership/ids/backend/pkg/config"
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
	hooks []func(old *Airport, new *Airport)
}

func NewAirportService(db *gorm.DB, c *cache.Cache) *AirportService {
	exp := time.Duration(config.GetConfig().Cache.DefaultExpiration.Airports) * time.Second

	return &AirportService{db: db, cache: c, exp: exp}
}

func (s *AirportService) AddHook(h func(old *Airport, new *Airport)) {
	s.hooks = append(s.hooks, h)
}

func (s *AirportService) callHooks(old *Airport, new *Airport) {
	for _, h := range s.hooks {
		h(old, new)
	}
}

func (s *AirportService) GetAirports() ([]*Airport, error) {
	var airports []*Airport
	if airport, err := s.cache.Get(cacheAll); err == nil {
		airports = airport.([]*Airport)
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

func (s *AirportService) GetAirport(id string) (*Airport, error) {
	var airport *Airport
	apt, err := s.cache.Get(cachePrefix + "/" + id)
	if err == nil {
		airport = apt.(*Airport)
	} else if err == cache.ErrorKeyNotFound {
		if err := s.db.Model(Airport{}).Where(Airport{FAAID: id}).Or(Airport{ICAOID: id}).First(&airport).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrInvalidAirport
			}
		}
		if err := s.cache.Set(cachePrefix+"/"+airport.FAAID, airport, s.exp); err != nil {
			return nil, err
		}
		if err := s.cache.Set(cachePrefix+"/"+airport.ICAOID, airport, s.exp); err != nil {
			return nil, err
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

	s.callHooks(nil, airport)

	return nil
}

func (s *AirportService) UpdateAirport(airport *Airport) error {
	var oldAirport *Airport
	oldAirport, err := s.GetAirport(airport.FAAID)
	if err != nil {
		return err
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

	s.callHooks(&airport, nil)

	return nil
}
