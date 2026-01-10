package main

import (
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"
)

func listenForMail() {

	// runs indefinitely and fire the message every time when triggered	in the background

	//we run go routine with anonymous function
	go func() {
		//infinitive loop
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()

}

// this function just send email
func sendMsg(m models.MailData) {
	//we'll install dummy server
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false //we close connection after every msg sent
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// we need a client
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	// we need to define an (empty)  message
	email := mail.NewMSG()
	// setting from, to and subject
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	//sending email
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}

}
