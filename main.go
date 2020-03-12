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

	errorMessage := "Cannot process the data"

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
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_, _ = b.Send(m.Sender, errorMessage)
			return
		}
		enc := base64.StdEncoding.EncodeToString(bts)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(1000)
		s := strconv.Itoa(n)
		file, err := os.Create(s + ".txt")
		if file != nil {
			_, err = file.Write([]byte(enc))
			log.Println("File", s+".txt", "has been created")
			err = file.Close()
		}
		p := &tb.Document{File: tb.FromDisk(s + ".txt"), FileName: "base64.txt"}
		_, _ = b.Send(m.Sender, p)
		err = os.Remove(s + ".txt")
		log.Println("File", s+".txt", "has been removed")
	})

	b.Handle(tb.OnDocument, func(m *tb.Message) {
		docId := m.Document.FileID
		url, err := b.FileURLByID(docId)
		if err != nil {
			_, _ = b.Send(m.Sender, errorMessage)
			return
		}
		log.Println(url)
		doc, err := http.Get(url)
		if err != nil {
			_, _ = b.Send(m.Sender, errorMessage)
			return
		}
		bts, err := ioutil.ReadAll(doc.Body)
		contentType := http.DetectContentType(bts)
		log.Println(contentType)
		if contentType != "text/plain; charset=utf-8" {
			_, _ = b.Send(m.Sender, errorMessage)
			return
		} else {
			dec, err := base64.StdEncoding.DecodeString(string(bts))
			if err != nil {
				_, _ = b.Send(m.Sender, errorMessage)
				return
			}
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(1000)
			s := strconv.Itoa(n)
			file, err := os.Create(s + ".jpg")
			if file != nil {
				_, err = file.Write(dec)
				log.Println("File", s+".jpg", "has been created")
				err = file.Close()
			}
			p := &tb.Photo{File: tb.FromDisk(s + ".jpg")}
			_, _ = b.Send(m.Sender, p)
			err = os.Remove(s + ".jpg")
			log.Println("File", s+".jpg", "has been removed")
		}
	})

	b.Start()
}
