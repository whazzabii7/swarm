package main

import (
	"fmt"
	"github.com/whazzabii7/swarm/internal/mf"
	"github.com/whazzabii7/swarm/internal/ui"
)

func main() {
	uiDone := make(chan bool)
	ui.InitUI(uiDone)

	ui.Log(ui.LevelInfo, "MAIN", `
   _____      S tructure.
  / ___/      W orkflow.
  \__ \       A utomation.
 ___/ /       R esilience.
/____/        M ainframe.
	`)
	ui.Log(ui.LevelInfo, "MAIN", ">>> Starting Swarm Mainframe...")

	done := make(chan bool)
	mfInstance := mf.NewMainframe()

	go mfInstance.Start(done)
	<-done
	ui.MainUI.Stop()
	<-uiDone
	fmt.Println(">>> successfully shut down.]")
}
