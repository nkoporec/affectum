package main

import (
	"time"

	"github.com/getlantern/systray"
	"github.com/nkoporec/affectum/icon"
	"github.com/nkoporec/affectum/utils"
	"github.com/skratchdot/open-golang/open"
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
	mFiles := systray.AddMenuItem("Files", "Saved files")
	mConfig := systray.AddMenuItem("Configuration", "Edit configuration")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	// Set up the dir if needed.
	utils.CreateDir()

	go func() {
		go executeScanMailJob()
		for {
			select {
			case <-mFiles.ClickedCh:
				err := open.Run(utils.GetAttachmentDir())
				if err != nil {
					utils.Logger(err.Error())
				}
			case <-mConfig.ClickedCh:
				err := open.Run(utils.GetDir())
				if err != nil {
					utils.Logger(err.Error())
				}
			case <-mQuitOrig.ClickedCh:
				utils.Logger("Requesting quit")
				systray.Quit()
				return
			}
		}
	}()
}
