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
	CREATE TABLE IF NOT EXISTS bots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT UNIQUE,
		path TEXT,
		type TEXT,
		status TEXT,
		meta TEXT,
		last_check DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		bot_alias TEXT,
		payload TEXT,
		status TEXT,
		priority INTEGER DEFAULT 0,
		dependency TEXT,
		retry_count INTEGER DEFAULT 0,
		max_retries INTEGER DEFAULT 3,
		timeout INTEGER DEFAULT 300,
		result TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT,
		origin TEXT,
		message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatalf("[-] Fehler beim Erstellen der Tabellen: %v", err)
	}

	log.Println("[+] Swarm Datenbank erfolgreich initialisiert.")
}
