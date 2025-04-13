package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/fx"
	"log/slog"
)

type Params struct {
	fx.In

	DB     *sql.DB
	Logger *slog.Logger
}

//go:embed schema/*.sql
var migrationFiles embed.FS

func RunMigrations(p Params) error {
	sourceDriver, err := iofs.New(migrationFiles, "schema")
	if err != nil {
		p.Logger.Error("failed to load migration files: ", err.Error())
		return fmt.Errorf("failed to initialize migrations source driver: %w", err)
	}

	dbDriver, err := postgres.WithInstance(p.DB, &postgres.Config{})
	if err != nil {
		p.Logger.Error("failed to initialize postgres driver: ", err.Error())
		return fmt.Errorf("failed to initialize postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "online_programming_school", dbDriver)
	if err != nil {
		p.Logger.Error("failed to initialize migrate instance: ", err.Error())
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		p.Logger.Error("failed to run migrations: ", err.Error())
		return fmt.Errorf("migration up failed: %w", err)
	}

	err = sourceDriver.Close()
	if err != nil {
		p.Logger.Error("failed to close migrations sourceDriver", "error", err)
		return fmt.Errorf("failed to close migrations sourceDriver: %w", err)
	}

	err = dbDriver.Close()
	if err != nil {
		p.Logger.Error("failed to close migrations dbDriver", "error", err)
		return fmt.Errorf("failed to close migrations dbDriver: %w", err)
	}

	p.Logger.Info("migrations done")
	return nil
}
