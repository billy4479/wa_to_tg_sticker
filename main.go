package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	maxSize uint64 = 20 * 1024 * 1024
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", start(bot))

	bot.Handle(tb.OnDocument, handleZip(bot))

	bot.Start()
}
