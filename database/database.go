package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"net/http"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var DBPool *pgxpool.Pool

func MigrateProcess() error {
	dbConn, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Config.Database.Host,
			config.Config.Database.Port,
			config.Config.Database.User,
			config.Config.Database.Password,
			config.Config.Database.Name,
			config.Config.Database.Sslmode,
		),
	)

	if err != nil {
		return err
	}

	defer func(sqlDb *sql.DB) {
		_ = sqlDb.Close()
	}(dbConn)

	migrateSrc := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(embedMigrations),
	}

	_, err = migrate.Exec(dbConn, "postgres", migrateSrc, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}

func ConnectDB() error {
	pgx, err := pgxpool.Connect(
		context.Background(),
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d",
			config.Config.Database.Host,
			config.Config.Database.Port,
			config.Config.Database.User,
			config.Config.Database.Password,
			config.Config.Database.Name,
			config.Config.Database.Sslmode,
			config.Config.Database.MaxConnect,
		),
	)

	if err == nil {
		DBPool = pgx
	}

	return err
}

func DisconnectDB() {
	DBPool.Close()
}
