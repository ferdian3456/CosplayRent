package main

import (
	"cosplayrent/app"
	"log"
)

func main() {
	_ = app.NewDB()
	log.Println("Here is main file")
}
