package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/user"

	"github.com/emersion/go-message/mail"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
	"github.com/nkoporec/affectum/utils"
)

func main() {
	usr, err := user.Current()

	fmt.Println("Loading config ...")

	config, err := utils.LoadConfig(usr.HomeDir)
	if err != nil {
		log.Fatal("Can't load config")
	}

	fmt.Println("Config loaded!")

	fmt.Println("Connecting to server ...")

	c, err := client.DialTLS(config.MailHost+":"+config.MailPort, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't connected to server, error is: %s", err))
	}

	fmt.Println("Connected")

	defer c.Logout()

	// Login
	fmt.Println("Logging in ...")
	err = c.Login(config.MailUsername, config.MailPassword)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't login to server, error is: %s", err))
	}
	log.Println("Logged in")

	// Select INBOX
	mbox, err := c.Select(config.MailFolder, false)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cant retrieve the mail folder, err was: %s", err))
	}

	// Get messages.
	// This only fetches the latest 4 messages.
	// @todo: How to do this more generic?
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 3
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	section := &imap.BodySectionName{}

	messages := make(chan *imap.Message, from)
	done := make(chan error, 1)

	// Fetch the message with go routines.
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	for msg := range messages {
		// Create a new mail reader
		mr, err := mail.CreateReader(msg.GetBody(section))
		if err != nil {
			log.Fatal(err)
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()

			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				fmt.Println(fmt.Sprintf("Saving attachment: %s", filename))
				b, _ := ioutil.ReadAll(p.Body)
				err := ioutil.WriteFile(filename, b, 0777)

				if err != nil {
					log.Println("Error while trying to save attachment: ", err)
				}

				fmt.Println(fmt.Sprintf("Attachment saved: %s", filename))
			}
		}
	}

}
