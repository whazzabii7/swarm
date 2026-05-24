package db

import (
	"context"
	"fmt"
	"time"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (g *Guardian) processSaveBlueprint(ctx context.Context, bp models.BotBlueprint) error {
	ui.Logf(ui.LevelInfo, "DB-Guardian", "SQL-Action: Saving Bot '%s' (%s)\n", bp.Alias, bp.Path)
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
		return fmt.Errorf("%w for %s: %v", ErrSQLUsert, bp.Alias, err)
	}
	ui.Logf(ui.LevelInfo, "DB-Guardian", "[+] Blueprint '%s' syncronised.", bp.Alias)
	return nil
}

func (g *Guardian) handleGetBlueprint(ctx context.Context, bpAlias string) (*models.BotBlueprint, error) {
	ui.Logf(ui.LevelInfo, "DB-Guardian", "SQL-Action: Get Bot '%s'\n", bpAlias)
	query := `
	SELECT alias, path, type, version, description, last_scan
    FROM bot_blueprints
    WHERE alias = ?
    LIMIT 1;`

	var bp models.BotBlueprint

	err := DB.QueryRowContext(ctx, query, bpAlias).Scan(
		&bp.Alias,
		&bp.Path,
		&bp.Type,
		&bp.Version,
		&bp.Description,
		&bp.LastScan,
	)
	if err != nil {
		return nil, fmt.Errorf("%w for %s: %v", ErrSQLUsert, bp.Alias, err)
	}
	ui.Logf(ui.LevelInfo, "DB-Guardian", "[+] Blueprint '%s' loaded.", bp.Alias)
	return &bp, nil
}

func (g *Guardian) handleCheckBlueprints(ctx context.Context, bps []models.BotBlueprint) ([]models.BotBlueprint, error) {
	return nil, nil
}
