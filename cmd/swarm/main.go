package main

import (
	"fmt"

	"github.com/whazzabii7/swarm/internal/mf"
)

func main() {
	fmt.Println(`
   _____      S tructure.
  / ___/      W orkflow.
  \__ \       A utomation.
 ___/ /       R esilience.
/____/        M ainframe.
	`)
	fmt.Println(">>> Starting Swarm Mainframe...")
	done := make(chan bool)
	mfInstance := mf.NewMainframe()

	go mfInstance.Start(done)
	<-done
	fmt.Println(">>> Swarm successfully shut down.")
}
