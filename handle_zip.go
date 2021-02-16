package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func handleZip(bot *tb.Bot) func(*tb.Message) {
	return func(m *tb.Message) {

		if m.Document == nil {
			bot.Send(m.Sender, "Invalid file")
			return
		}

		if m.Document.MIME != "application/zip" {
			bot.Send(m.Sender, "This doesn't look like a zip...")
		}

		if uint64(m.Document.FileSize) > maxSize {
			bot.Send(m.Sender, fmt.Sprintf("File too big: max zip size %d", maxSize))
			return
		}

		dest, err := unzip(&m.Document.File, bot)
		if err != nil {
			bot.Send(m.Sender, fmt.Sprintf("An error occurred: %s"), err.Error())
			return
		}
	}
}
