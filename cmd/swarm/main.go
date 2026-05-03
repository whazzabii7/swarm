package main

import (
	"fmt"
	"log"
	"os"
	"encoding/json"

	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/bot"
)

func main() {
	fmt.Println(`
   _____      S tructure.
  / ___/      W orkflow.
  \__ \       A utomation.
 ___/ /       R esilience.
/____/        M ainframe.
	`)
	fmt.Println(">>> Swarm Mainframe wird gestartet...")

	// Sicherstellen, dass das data-Verzeichnis existiert
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	// DB Initialisieren
	db.InitDB("data/swarm.db")

	// Hier kommt später der Loop für den Event-Listener und Scheduler rein
	log.Println("[!] Mainframe läuft. Warte auf Befehle...")

	botManager := bot.NewBotManager()
	
	bps, err := botManager.SyncBlueprints()
	if err != nil {
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
		return
	}
	for _, bp := range bps {
		prettyJSON, _ := json.MarshalIndent(bp, "", "  ")
		fmt.Println(string(prettyJSON))
	}

	// Verhindert, dass das Programm sofort beendet (vorerst)
	select {} 
}
