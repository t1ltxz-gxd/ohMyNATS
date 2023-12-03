package redis

import (
	"github.com/go-redis/redis"
	c "ohMyNATS/pkg/client"
)

type Redis struct {
	c.Builder
	Addr     string
	Password string
	DB       int
}

func (r *Redis) New() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	_, err := client.Ping().Result()
	return client, err
}
