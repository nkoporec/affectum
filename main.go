package main

import (
	"time"

	"github.com/getlantern/systray"
	"github.com/nkoporec/affectum/icon"
	"github.com/nkoporec/affectum/utils"
)

func main() {
	// Start the systray.
	onExit := func() {
		utils.Logger("Affectum stopped!")
	}

	systray.Run(onReady, onExit)
}

func executeScanMailJob() {
	for {
		utils.ScanMail()
		time.Sleep(60 * time.Second)
	}
}

func onReady() {
	systray.SetIcon(icon.Base.Data)
	systray.SetTitle("")
	systray.SetTooltip("Affectum")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	// Set up the dir if needed.
	utils.CreateDir()

	go func() {
		go executeScanMailJob()

		<-mQuitOrig.ClickedCh
		utils.Logger("Requesting quit")
		systray.Quit()
	}()
}
