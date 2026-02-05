package main

import (
	"embed"

	"github.com/gclm/gclm-flow/gclm-engine/internal/assets"
)

// AssetsFS contains all embedded assets for the application
//
// Currently embedded:
//   - migrations/*.sql            Database migration files
//   - workflows/*.yaml            Default workflow definitions
//   - gclm_engine_config.yaml     Default configuration
//
// Future expansion (add when directories exist):
//   - web/static/*         Static files for web UI
//   - web/templates/*       HTML templates
//
//go:embed migrations/*.sql workflows/*.yaml gclm_engine_config.yaml
var AssetsFS embed.FS

// MigrationsFS returns the embedded migrations filesystem
// For use with goose: pass AssetsFS directly and use "migrations" as the path
func MigrationsFS() *embed.FS {
	return &AssetsFS
}

// GetEmbedFS returns the embedded filesystem for internal use
func GetEmbedFS() *embed.FS {
	return &AssetsFS
}

func init() {
	// Initialize the assets package with the filesystem
	assets.Init(AssetsFS)
}
