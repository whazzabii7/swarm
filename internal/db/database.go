package db

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite" // CGO-freier Treiber
)

var DB *sql.DB

// InitDB initialisiert die SQLite Datenbank und legt Tabellen an
func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("[-] Fehler beim Öffnen der Datenbank: %v", err)
	}

	// Tabellen-Schema
	const schema = `
	CREATE TABLE IF NOT EXISTS bot_blueprints (
		alias TEXT PRIMARY KEY,
		path TEXT NOT NULL,
		type TEXT,
		last_scan DATETIME
	);

	CREATE TABLE IF NOT EXISTS bot_instances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT,
		pid INTEGER,
		status TEXT,
		last_seen DATETIME,
		FOREIGN KEY(alias) REFERENCES bot_blueprints(alias)
	);

	CREATE TABLE IF NOT EXISTS task_ledger (
		id TEXT PRIMARY KEY,
		target_bot TEXT,
		payload TEXT,
		status TEXT,
		created_at DATETIME,
		FOREIGN KEY(target_bot) REFERENCES bot_blueprints(alias)
	);`

	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatalf("[-] Fehler beim Erstellen der Tabellen: %v", err)
	}

	log.Println("[+] Swarm Datenbank erfolgreich initialisiert.")
}
