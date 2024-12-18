package config

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/knadh/koanf/v2"
)

func NewMemcacheClient(config *koanf.Koanf) *memcache.Client {
	port := config.String("cache.server_port")
	client := memcache.New(":" + port)

	return client
}
