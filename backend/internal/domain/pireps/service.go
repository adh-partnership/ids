package pireps

import (
	"errors"
	"fmt"
	"time"

	"github.com/adh-partnership/ids/backend/internal/cache"
	"github.com/adh-partnership/ids/backend/pkg/config"
	"gorm.io/gorm"
)

const (
	cachePrefix = "/pireps"
	cacheAll    = cachePrefix + "/_all"
)

var (
	ErrInvalidPIREP = errors.New("invalid airport")
)

type PIREPService struct {
	db    *gorm.DB
	cache *cache.Cache
	exp   time.Duration
	hooks []func(new *PIREP)
}

func NewPIREPService(db *gorm.DB, c *cache.Cache) *PIREPService {
	exp := time.Duration(config.GetConfig().Cache.DefaultExpiration.PIREPs) * time.Second

	return &PIREPService{db: db, cache: c, exp: exp}
}

func (s *PIREPService) AddHook(h func(new *PIREP)) {
	s.hooks = append(s.hooks, h)
}

func (s *PIREPService) callHooks(new *PIREP) {
	for _, h := range s.hooks {
		h(new)
	}
}

func (s *PIREPService) GetPIREPs() ([]*PIREP, error) {
	var pireps []*PIREP
	if pirep, err := s.cache.Get(cacheAll); err == nil {
		pirep = pirep.([]*PIREP)
	} else if err == cache.ErrorKeyNotFound {
		if err := s.db.Find(&pireps).Error; err != nil {
			return nil, err
		}
		if err := s.cache.Set(cacheAll, pireps, s.exp); err != nil {
			return nil, err
		}
	}

	return pireps, nil
}

func (s *PIREPService) GetPIREP(id int64) (*PIREP, error) {
	var pirep *PIREP
	pi, err := s.cache.Get(fmt.Sprintf("%s/%d", cachePrefix, id))
	if err == nil {
		pirep = pi.(*PIREP)
	} else if err == cache.ErrorKeyNotFound {
		if err := s.db.Model(PIREP{}).Where(PIREP{ID: id}).First(&pirep).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrInvalidPIREP
			}
		}
		if err := s.cache.Set(s.getKey(id), pirep, s.exp); err != nil {
			return nil, err
		}
	}

	return pirep, nil
}

func (s *PIREPService) CreatePIREP(pirep *PIREP) error {
	pirep.Raw = pirep.ToString()
	if err := s.db.Create(&pirep).Error; err != nil {
		return err
	}
	if err := s.cache.Set(s.getKey(pirep.ID), pirep, s.exp); err != nil {
		return err
	}
	if err := s.cache.Delete(cacheAll); err != nil {
		return err
	}

	s.callHooks(pirep)

	return nil
}

func (s *PIREPService) DeletePIREP(id int64) error {
	var pirep *PIREP
	if err := s.db.Where(PIREP{ID: id}).First(&pirep).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&pirep).Error; err != nil {
		return err
	}
	if err := s.cache.Delete(s.getKey(id)); err != nil {
		return err
	}
	if err := s.cache.Delete(cacheAll); err != nil {
		return err
	}

	return nil
}

func (s *PIREPService) DeleteExpiredPIREPs() error {
	var pireps []*PIREP
	// Find pireps more than an hour old
	if err := s.db.Where("tm < ?", time.Now().Add(-1*time.Hour)).Find(&pireps).Error; err != nil {
		return err
	}
	if len(pireps) > 0 {
		if err := s.db.Delete(&pireps).Error; err != nil {
			return err
		}
		for _, pirep := range pireps {
			if err := s.cache.Delete(s.getKey(pirep.ID)); err != nil {
				if !errors.Is(err, cache.ErrorKeyNotFound) {
					return err
				}
			}
		}
		if err := s.cache.Delete(cacheAll); err != nil {
			return err
		}
	}

	return nil
}

func (s *PIREPService) getKey(id int64) string {
	return fmt.Sprintf("%s/%d", cachePrefix, id)
}
