package main

import (
	"fmt"

	"time"

	"github.com/nkoporec/affectum/utils"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func main() {
	// Start the systray.
	onExit := func() {
		fmt.Println("Stopped!")
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

	go func() {
		go executeScanMailJob()

		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
	}()
}
