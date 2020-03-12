package main

import (
	"encoding/base64"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
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

	errorMessage := "Cannot process the Image"

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
			_, _ = b.Send(m.Sender, errorMessage)
			return
		}
		log.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			_, _ = b.Send(m.Sender, errorMessage)
			return
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_, _ = b.Send(m.Sender, errorMessage)
			return
		}
		enc := base64.StdEncoding.EncodeToString(bytes)
		log.Println(enc)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(1000)
		s := strconv.Itoa(n)
		file, err := os.Create(s + ".txt")
		if file != nil {
			_, err = file.Write([]byte(enc))
			err = file.Close()
		}
		p := &tb.Document{File: tb.FromDisk(s + ".txt"), FileName: "base64.txt"}
		_, _ = b.Send(m.Sender, p)
		err = os.Remove(s + ".txt")
	})

	b.Start()
}
