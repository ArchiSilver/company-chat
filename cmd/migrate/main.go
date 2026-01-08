package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"

	"company-chat/internal/config"
	"company-chat/internal/logging"
)

func main() {
	cfg := config.LoadConfig()
	// init logger for migrate tool
	logging.Init(false)
	defer logging.Sync()

	db, err := sql.Open("pgx", cfg.GetDBConnectionString())
	if err != nil {
		logging.L.Fatalf("open db: %v", err)
	}
	defer db.Close()

	migDir := "migrations"
	entries := []string{}
	err = filepath.WalkDir(migDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sql" && filepath.Base(path)[len(filepath.Base(path))-7:] == ".up.sql" {
			entries = append(entries, path)
		}
		return nil
	})
	if err != nil {
		logging.L.Fatalf("walk migrations: %v", err)
	}

	for _, f := range entries {
		logging.L.Infof("applying migration: %s", f)
		b, err := os.ReadFile(f)
		if err != nil {
			logging.L.Fatalf("read %s: %v", f, err)
		}
		if _, err := db.ExecContext(context.Background(), string(b)); err != nil {
			logging.L.Fatalf("exec migration %s: %v", f, err)
		}
	}
	fmt.Println("✅ Миграции применены")
}
