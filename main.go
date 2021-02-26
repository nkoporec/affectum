package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"time"

	"github.com/nkoporec/affectum/utils"
	"github.com/sevlyar/go-daemon"
)

var (
	signal = flag.String("s", "", `Send signal to the daemon:
  stop — fast shutdown`)
)

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func main() {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)

	cntxt := &daemon.Context{
		PidFileName: "affectum.pid",
		PidFilePerm: 0644,
		LogFileName: "affectum.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon affectum]"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	fmt.Println("Affectum starting ...")
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}

	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("affectum started")

	go executeScanMailJob()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	log.Println("affectum terminated")
}

func executeScanMailJob() {
LOOP:
	for {
		utils.ScanMail()
		time.Sleep(60 * time.Second)
		select {
		case <-stop:
			break LOOP
		default:
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")

	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}

	return daemon.ErrStop
}
