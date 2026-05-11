package db

import (
    "context"
	"fmt"
	"log"
	"time"

	"github.com/whazzabii7/swarm/internal/models"
)

func (g *Guardian) processSaveBlueprint(ctx context.Context, bp models.BotBlueprint) {
	fmt.Printf("[DB-Guardian] SQL-Action: Saving Bot '%s' (%s)\n", bp.Alias, bp.Path)
	query := `
	INSERT INTO bot_blueprints (alias, path, type, version, description, last_scan)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(alias) DO UPDATE SET
			path = excluded.path,
			type = excluded.type,
			version = excluded.version,
			description = excluded.description,
			last_scan = excluded.last_scan;`

	_, err := DB.ExecContext(ctx, query, bp.Alias, bp.Path, bp.Type, bp.Version, bp.Description, bp.LastScan.Format(time.RFC3339))
	if err != nil {
		log.Printf("[-] Guardian Error with SQL-Upsert for %s: %v", bp.Alias, err)
		return
	}
	log.Printf("[+] Guardian Blueprint '%s' syncronised.", bp.Alias)
}
