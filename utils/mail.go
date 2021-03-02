package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"github.com/emersion/go-message/mail"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
)

func ScanMail() bool {
	usr, err := user.Current()

	Logger("Loading config ...")

	config, err := LoadConfig(usr.HomeDir)
	if err != nil {
		Logger("Can't load config")
		log.Fatal("")
	}

	Logger("Config loaded!")

	Logger("Loading database ...")
	db := SetupDatabase()
	if db != true {
		Logger("Loading database failed!")
		log.Fatal("")
	}

	Logger("Database loaded!")

	Logger("Connecting to server ...")

	c, err := client.DialTLS(config.MailHost+":"+config.MailPort, nil)
	if err != nil {
		Logger(fmt.Sprintf("Can't connected to server, error is: %s", err))
		log.Fatal()
	}

	Logger("Connected")

	defer c.Logout()

	// Login
	Logger("Logging in ...")
	err = c.Login(config.MailUsername, config.MailPassword)
	if err != nil {
		Logger(fmt.Sprintf("Can't login to server, error is: %s", err))
	}
	Logger("Logged in!")

	// Select INBOX
	mbox, err := c.Select(config.MailFolder, false)
	if err != nil {
		Logger(fmt.Sprintf("Cant retrieve the mail folder, err was: %s", err))
	}

	Logger("Scanning folder ...")

	// Get ALL messages.
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)
	section := &imap.BodySectionName{}

	messages := make(chan *imap.Message, mbox.Messages)
	done := make(chan error, 1)

	// Fetch the message with go routines.
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem(), imap.FetchUid}, messages)
	}()

	for msg := range messages {
		// Don't process old mails.
		if MailExists(config.MailFolder, msg.Uid) == true {
			continue
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(msg.GetBody(section))
		if err != nil {
			Logger(err.Error())
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()

			if err == io.EOF {
				break
			} else if err != nil {
				Logger(err.Error())
			}

			switch h := p.Header.(type) {
			case *mail.AttachmentHeader:
				// log the message in database.
				InsertMail(config.MailFolder, msg.Uid)

				filename, _ := h.Filename()
				Logger(fmt.Sprintf("Saving attachment: %s", filename))
				b, _ := ioutil.ReadAll(p.Body)

				attachmentFolder := config.AttachmentFolderPath
				if attachmentFolder == "" {
					attachmentFolder = filepath.Join(usr.HomeDir, "/affectum/files")
				}

				attachment := filepath.Join(attachmentFolder, filename)
				err := ioutil.WriteFile(attachment, b, 0777)

				if err != nil {
					Logger(err.Error())
				}

				Logger(fmt.Sprintf("Attachment saved: %s", filename))
			}
		}
	}

	Logger("Scan completed!")

	return true
}
