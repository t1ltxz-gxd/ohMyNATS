package main

import (
	"ohMyNATS/internal/app"
	"ohMyNATS/internal/config"
)

func main() {
	cfg := config.InitConfig()
	app.Serve(cfg.Env, cfg.BackEndPort)
}
