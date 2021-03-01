package main

import (
	"time"
	"github.com/nkoporec/affectum/utils"
	"github.com/getlantern/systray"
	"github.com/nkoporec/affectum/icon"
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
	systray.SetTemplateIcon(icon.Data, icon.Data)
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
