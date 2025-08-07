package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"voteflix/api/internal/app"
)

func main() {
	cfg := app.Init().Config()
	m, connectErr := migrate.New(
		"file://migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=public",
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresDb,
		),
	)

	if connectErr != nil {
		log.Fatal(connectErr)
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("No change made by migration scripts")
		} else {
			log.Fatal(err)
		}
	}
}
