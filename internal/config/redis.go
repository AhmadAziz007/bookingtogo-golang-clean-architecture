package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper) *redis.Client {
	host := viper.GetString("redis.host")
	database := viper.GetInt("redis.database")

	var client = redis.NewClient(&redis.Options{
		Addr: host,
		DB:   database,
	})
	return client
}
