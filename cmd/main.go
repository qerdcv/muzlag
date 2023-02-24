package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qerdcv/muzlag/internal/bot"
)

func main() {
	b, err := bot.New()
	if err != nil {
		log.Fatalln(err)
	}

	if err = b.Open(); err != nil {
		log.Fatalln(err)
	}

	if err = b.RegisterCommands(); err != nil {
		log.Fatalln(err)
	}

	log.Println("Bot is started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err = b.Close(); err != nil {
		log.Fatalln(err)
	}
}
