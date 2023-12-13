package redis

import (
	"fmt"

	"github.com/adh-partnership/ids/backend/pkg/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func New(rc *config.Cache) *Redis {
	if rc.Driver == "redis" {
		r := redis.NewClient(&redis.Options{
			Addr:     rc.Host + ":" + fmt.Sprint(rc.Port),
			Password: rc.Password,
			DB:       rc.DB,
		})
		return &Redis{r}
	}

	return nil
}
