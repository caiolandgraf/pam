package main

import (
	"log"

	"github.com/eduardofuncao/pam/internal/config"
)

func main() {
	cfg, err := config.LoadConfig(config.CfgFile)
	if err != nil {
		log.Fatal("Could not load config file", err)
	}

	app := NewApp(cfg)
	app.Run()
}
