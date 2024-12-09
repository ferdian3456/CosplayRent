package app

import (
	"cosplayrent/helper"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
	"os"
)

func NewClient() *memcache.Client {
	var err error = godotenv.Load("../.env")
	helper.PanicIfError(err)

	MEMCACHED_PORT := os.Getenv("MEMCACHED_PORT")
	client := memcache.New(MEMCACHED_PORT)

	return client
}
