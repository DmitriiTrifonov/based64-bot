package main

import (
	"encoding/base64"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	publicURL := os.Getenv("PUBLIC_URL")
	token := os.Getenv("TOKEN")

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	pref := tb.Settings{
		Token:  token,
		Poller: webhook,
	}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	message := "This is the based64 bot.\n" +
		"It converts Images to Base64 strings and vice versa\n" +
		"Send a photo to start"

	b.Handle("/start", func(m *tb.Message) {
		_, _ = b.Send(m.Sender, message)
	})

	b.Handle("/help", func(m *tb.Message) {
		_, _ = b.Send(m.Sender, message)
	})

	b.Handle("/about", func(fm *tb.Message) {
		_, _ = b.Send(fm.Sender, message)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		photoUrl := m.Photo.File.FileID
		url, err := b.FileURLByID(photoUrl)
		if err != nil {
			_, _ = b.Send(m.Sender, "Cannot process the photo")
			return
		}
		log.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			_, _ = b.Send(m.Sender, "Cannot process the photo")
			return
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_, _ = b.Send(m.Sender, "Cannot process the photo")
			return
		}
		enc := base64.StdEncoding.EncodeToString(bytes)
		log.Println(enc)
		_, _ = b.Send(m.Sender, enc)
	})

	b.Start()
}
