package app

import "github.com/bradfitz/gomemcache/memcache"

func NewClient() *memcache.Client {
	client := memcache.New(":11211")

	return client
}
