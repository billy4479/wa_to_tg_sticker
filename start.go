package main

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func start(bot *tb.Bot) func(*tb.Message) {
	return func(m *tb.Message) {
		_, err := bot.Send(m.Sender, `
Welcome! This bot will convert for you your Whatsapp stickers into a Telegram Sticker Pack.
		
1) Create a group chat on Whatsapp with just yourself in it
2) Send all the stickers you want to convert to that group
3) Export the chat of that group
4) Make a zip
5) Send it to me
6) I'll ask you a name and a default emoji (you can change them later)
7) Done!
		`)
		if err != nil {
			log.Fatal(err)
		}
	}
}
