package database

import (
	"database/sql"
	"embed"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"io/fs"
	"net/http"
)

//go:embed migrations/*.sql
var migrations embed.FS

func DBMigrate(
	host string,
	port uint16,
	user string,
	password string,
	dbname string,
	sslmode string,
) (int, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host,
			port,
			user,
			password,
			dbname,
			sslmode,
		),
	)

	if err != nil {
		return 0, err
	}

	defer func() {
		_ = db.Close()
	}()

	migrateSrc := &migrate.HttpFileSystemMigrationSource{
		FileSystem: func() http.FileSystem {
			dist, err := fs.Sub(migrations, "migrations")
			if err != nil {
				panic(err)
			}
			return http.FS(dist)
		}(),
	}

	applied, err := migrate.Exec(db, "postgres", migrateSrc, migrate.Up)
	if err != nil {
		return 0, err
	}

	return applied, err
}
