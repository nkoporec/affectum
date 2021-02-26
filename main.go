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

	fmt.Println("Loading database ...")
	db := utils.SetupDatabase()
	if db != true {
		log.Fatal("Loading database failed!")
	}

	fmt.Println("Database loaded!")

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

	// Get ALL messages.
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)
	section := &imap.BodySectionName{}

	messages := make(chan *imap.Message, mbox.Messages)
	done := make(chan error, 1)

	// Fetch the message with go routines.
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	// @todo: Grab each message id (msg.Uid) save it to some database and then
	// when we loop through messages make sure that we only process each message
	// once.
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
				// log the message in database.
				utils.InsertMail(config.MailFolder, msg.Uid)

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
