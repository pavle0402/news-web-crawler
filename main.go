package main

import (
	"crawler/api"
	"crawler/config"
	"crawler/db"
)

func main() {
	cfg := config.LoadConfig()
	db.Connect(cfg)

	api.StartServer()
}
