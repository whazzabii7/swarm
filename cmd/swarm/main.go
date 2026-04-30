package main

import (
	"fmt"
	"log"
	"os"

	"github.com/whazzabii7/swarm/internal/db"
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
	
	// Verhindert, dass das Programm sofort beendet (vorerst)
	select {} 
}
