package main

import (
	"flag"
	tgClient "github.com/sungatuly22/readAdvisorBot/clients/telegram"
	"github.com/sungatuly22/readAdvisorBot/consumer/eventConsumer"
	"github.com/sungatuly22/readAdvisorBot/events/telegram"
	"github.com/sungatuly22/readAdvisorBot/storage/files"
	"log"
)

const (
	tgHostBot   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(tgClient.New(tgHostBot, token()), files.New(storagePath))
	log.Print("service started")
	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func token() string {
	t := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()
	if *t == "" {
		log.Fatal("token is empty")
	}
	return *t
}
