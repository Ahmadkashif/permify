package factories

import (
	"github.com/Permify/permify/internal/repositories/postgres/migrations"
	"github.com/Permify/permify/pkg/database"
)

// MigrationFactory - Get migrations according to given engine
func MigrationFactory(engine database.Engine) []string {
	switch engine {
	case database.POSTGRES:
		return migrations.GetMigrations()
	case database.MEMORY:
		return []string{}
	default:
		return []string{}
	}
}
