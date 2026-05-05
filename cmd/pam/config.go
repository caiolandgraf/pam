package main

import (
	"fmt"
	"log"

	"github.com/caiolandgraf/pam/internal/configui"
	"github.com/caiolandgraf/pam/internal/styles"
)

func (a *App) handleConfig() {
	saved, err := configui.Run(a.config)
	if err != nil {
		log.Fatalf("Config UI error: %v", err)
	}
	if saved {
		fmt.Println(styles.Success.Render("✓ Configuration saved successfully"))
	}
}
