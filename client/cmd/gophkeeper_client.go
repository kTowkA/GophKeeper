package main

import (
	"context"
	"log"

	"github.com/kTowkA/GophKeeper/client"
	"github.com/kTowkA/GophKeeper/client/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Run(context.Background(), *cfg)
	if err != nil {
		log.Fatal(err)
	}
}
