package main

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
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

		filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".webp") {
				return nil
			}

			r, err := os.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()

			wa, err := webp.Decode(r)
			if err != nil {
				return err
			}

			buf := bytes.NewBuffer(make([]byte, int(info.Size())))
			err = png.Encode(buf, wa)
			if err != nil {
				return err
			}

			stk := bot.UploadStickerFile(m.Sender, &tb.FromReader(buf))

			return nil
		})
	}
}
