package config

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/knadh/koanf/v2"
)

func NewMemcacheClient(config *koanf.Koanf) *memcache.Client {
	port := config.String("MEMCACHED_SERVER_PORT")
	client := memcache.New(":" + port)

	return client
}
