package app

import (
	"cosplayrent/helper"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func NewClient() *memcache.Client {
	var err error = godotenv.Load("../.env")
	helper.PanicIfError(err)

	MEMCACHED_SERVER_PORT := os.Getenv("MEMCACHED_SERVER_PORT")
	log.Println(MEMCACHED_SERVER_PORT)
	client := memcache.New(":" + MEMCACHED_SERVER_PORT)

	return client
}
