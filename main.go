package main

import (
	"log"

	"github.com/qerdcv/muzlag/bot"
	"github.com/qerdcv/muzlag/config"
)

func main() {
	cfg, err := config.New()

	if err != nil {
		log.Fatal(err.Error())
	}

	bot.Start(cfg)

	<-make(chan struct{})
}
