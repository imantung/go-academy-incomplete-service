package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

const (
	DbName     = "incomplete_service_development"
	DbUser     = "user"
	DbPassword = "password"
	DbHost     = "localhost"
	DbPort     = "5432"
)

func main() {
	app := cli.App{
		Name:        "incomplete-service",
		Description: "hello-world service",
		Version:     "0.1",
		Commands: []cli.Command{
			{
				Name:      "start",
				ShortName: "s",
				Usage:     "start HTTP Server",
				Action:    actionStart,
			},
			{
				Name:      "migrate up",
				ShortName: "u",
				Usage:     "up migration",
				Action:    actionMigrationUp,
			},
			{
				Name:      "migrate down",
				ShortName: "d",
				Usage:     "down migration",
				Action:    actionMigrationDown,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Some error occurred: %s", err.Error())
	}
}

func actionStart(c *cli.Context) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := map[string]string{
			"hello": "world",
		}

		b, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("err.Error()"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func actionMigrationUp(c *cli.Context) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DbUser, DbPassword, DbHost, DbPort, DbName)

	db, err := sql.Open("postgres", connString)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	fatalIfError(err)

	err = m.Up()
	fatalIfError(err)
}

func actionMigrationDown(c *cli.Context) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DbUser, DbPassword, DbHost, DbPort, DbName)

	db, err := sql.Open("postgres", connString)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	fatalIfError(err)

	err = m.Down()
	fatalIfError(err)
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
